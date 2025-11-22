package services_v1

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repositories"
	"user-management-api/internal/services"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	// Add necessary fields, e.g., database connection, config, etc.
	userRepo     repositories.UserRepo
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
	mailService  mail.EmailProviderService
}
type LoginAttempt struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	mutex           = &sync.Mutex{}
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute // time to live for login attempts
	MaxLoginAttempt = 5               // max 5 attempts in the TTL
)

func NewAuthService(userRepo repositories.UserRepo, tokenService auth.TokenService, cache cache.RedisCacheService, mailService mail.EmailProviderService) services.AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
		cacheService: cache,
		mailService:  mailService,
	}
}

func (as *authService) getClientIP(c *gin.Context) string {
	// implement logic to get client IP
	ip := c.ClientIP()
	if ip == "" {
		ip = c.Request.RemoteAddr // lay ip remote if client ip is empty
	}
	return ip
}

func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()
	client, exist := clients[ip]
	if !exist {
		// get env variable for rate limiter
		// reqSec := utils.GetIntEnv("RATE_LIMITER_REQUEST_SEC", 5)
		// brustSec := utils.GetIntEnv("RATE_LIMITER_REQUEST_BRUST", 10)

		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt) // 5 requests per 5 minutes
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		//log.Printf("A client with IP %s - {limiter: %+v , lastseen: %+s}", ip, newClient.Limiter, newClient.LastSeen)
		return limiter
	}

	//log.Printf("A client with IP %s - {limiter: %+v , lastseen: %+s}", ip, client.Limiter, client.LastSeen)

	client.LastSeen = time.Now()
	return client.Limiter
}

func (as *authService) checkLoginAttempt(ip string) error {
	// get client IP
	limiter := as.getLoginAttempt(ip)

	if !limiter.Allow() {
		// exceeded rate limit
		return utils.NewError("Too many login attempts. Pls try again later", utils.ErrorTooManyRequests)
	}
	return nil
}

func (as *authService) CleanUpClients(ip string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(clients, ip)
}

// Implement authentication-related methods here, e.g., Login, Logout, etc.
func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	// Implementation for login
	context := ctx.Request.Context()

	ip := as.getClientIP(ctx)
	// check login attempt
	if err := as.checkLoginAttempt(ip); err != nil {
		return "", "", 0, err
	}
	email = utils.NormalizeString(email)
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip) // record login attempt
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		as.getLoginAttempt(ip) // record login attempt
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.NewWrapError("Unable to create access token", utils.ErrorCodeInternal, err)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.NewWrapError("Unable to create refresh token", utils.ErrorCodeInternal, err)
	}

	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.NewWrapError("Cannot save refresh token", utils.ErrorCodeInternal, err)
	}
	// cleanup record login attempt when login success
	as.CleanUpClients(ip)
	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	// Implementation for logout

	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {

		return utils.NewError("Authorization header missing or invalid", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid access token", utils.ErrCodeUnauthorized)
	}
	// Blacklist the access token
	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cacheService.Set(key, "revoke", ttl)
	}

	// validate the refresh token
	_, err = as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.NewWrapError("Invalid refresh token ", utils.ErrCodeUnauthorized, err)
	}
	// Revoke the old refresh token
	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return utils.NewWrapError("Cannot revoke old refresh token", utils.ErrorCodeInternal, err)
	}
	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshToken string) (string, string, int, error) {
	// Implementation for refresh token
	context := ctx.Request.Context()
	// Validate the refresh token
	token, err := as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", 0, utils.NewWrapError("Invalid refresh token or revoked", utils.ErrCodeUnauthorized, err)
	}
	// Get the user associated with the refresh token
	userUuid, _ := uuid.Parse(token.UserUUID)
	user, err := as.userRepo.GetUserByUUID(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewWrapError("User not found", utils.ErrCodeUnauthorized, err)
	}
	// Generate new access and refresh tokens
	newAccessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.NewWrapError("Unable to create access token", utils.ErrorCodeInternal, err)
	}
	// Generate new refresh token
	newRefreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.NewWrapError("Unable to create refresh token", utils.ErrorCodeInternal, err)
	}

	// Revoke the old refresh token
	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.NewWrapError("Cannot revoke old refresh token", utils.ErrorCodeInternal, err)
	}
	// Store the new refresh token and invalidate the old one
	if err := as.tokenService.StoreRefreshToken(newRefreshToken); err != nil {
		return "", "", 0, utils.NewWrapError("Cannot save refresh token", utils.ErrorCodeInternal, err)
	}

	return newAccessToken, newRefreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) RequestForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()
	rateLimitKey := fmt.Sprintf("reset:ratelimit:%s", email)
	// check exist reset key in cache
	exist, err := as.cacheService.Exists(rateLimitKey)
	if err == nil && exist {
		return utils.NewError("Please wait before requesting another password reset.", utils.ErrorTooManyRequests)
	}

	email = utils.NormalizeString(email)

	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return utils.NewError("Invalid email or password", utils.ErrorCodeNotFound)
	}
	// send randome code to email
	str, err := utils.GenerateRandomeString(16)
	if err != nil {
		return utils.NewError("Failed to  generate randome string", utils.ErrorCodeInternal)
	}

	// save/set life string in cache
	err = as.cacheService.Set("reset:"+str, user.UserUuid, 1*time.Hour)
	if err != nil {
		return utils.NewError("Failed to store reset string password", utils.ErrorCodeInternal)
	}

	err = as.cacheService.Set(rateLimitKey, "1", 5*time.Minute)
	if err != nil {
		return utils.NewError("Failed to store ratelimit forget password", utils.ErrorCodeInternal)
	}

	resetLink := fmt.Sprintf("http://fontend.domain/reset-password?token=%s", str)
	logger.Log.Info().Msg(resetLink)
	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: email},
		},
		Subject: "Password Reset Request",
		Text: fmt.Sprintf("Hi %s, \n\n You requested to reset your password. Please click the link below to reset it:\n%s\n\n The link will expire in 1 hour. \n\n Best regard, \nCode With DatDev Team",
			user.UserEmail,
			resetLink),
	}
	///
	if err := as.mailService.SendMail(context, mailContent); err != nil {
		return utils.NewWrapError("Failed to send mail", utils.ErrorCodeInternal, err)
	}

	return nil
}

func (as *authService) ResetPassword(ctx *gin.Context, token string, newPassword string) error {
	context := ctx.Request.Context()
	var userUUID string
	err := as.cacheService.Get("reset:"+token, &userUUID)
	//log.Printf("---- userUUID:%s", userUUID)
	if err == redis.Nil || userUUID == "" {
		return utils.NewError("Invalid or expired token", utils.ErrorCodeNotFound)
	}
	if err != nil {
		return utils.NewError("Failed to get reset token", utils.ErrorCodeInternal)
	}
	useruuid, err := uuid.Parse(userUUID)
	if err != nil {
		return utils.NewError("UUid is invalid", utils.ErrorCodeInternal)
	}
	// new Pass
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.NewError("Failed to hash password", utils.ErrorCodeInternal)
	}
	var input = sqlc.UpdatePasswordParams{
		UserPassword: string(hashed),
		UserUuid:     useruuid,
	}
	_, err = as.userRepo.UpdatePassword(context, input)
	if err != nil {
		return utils.NewError("Update password failed", utils.ErrorCodeInternal)
	}
	// del cache reset:
	as.cacheService.Clear("reset:" + token)
	return nil
}
