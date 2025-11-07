package services_v1

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repositories"
	"user-management-api/internal/services"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	userRepo repositories.UserRepo
	cache    cache.RedisCacheService
}

func NewUsersService(repo repositories.UserRepo, redis *redis.Client) services.UsersService {
	return &usersService{
		userRepo: repo,
		cache:    cache.NewRedisCacheService(redis),
	}
}

func (us *usersService) GetUsers(ctx *gin.Context, search *string, order_by, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	if sort == "" {
		sort = "desc"
	}
	if order_by == "" {
		order_by = "user_created_at"
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		envLimit := utils.GetIntEnv("LIMIT_ITEM_PER_PAGE", 10)
		limit = int32(envLimit)
	}

	///caculate offset
	offset := (page - 1) * limit
	//check cache data
	valueSearch := reflect.ValueOf(search).String()
	cacheKey := us.generateCacheKey(valueSearch, order_by, sort, page, limit, deleted)
	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}
	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Users != nil {
		log.Printf("write to cache key: %s", cacheKey)
		return cacheData.Users, cacheData.Total, nil
	}

	users, err := us.userRepo.GetUsersV2(context, search, order_by, sort, offset, limit, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.NewWrapError("Failed to get users", utils.ErrorCodeInternal, err)
	}
	total, err := us.userRepo.UsersCount(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.NewWrapError("Failed to count users", utils.ErrorCodeInternal, err)
	}

	// write to cache
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: users,
		Total: int32(total),
	}
	us.cache.Set(cacheKey, cacheData, 3*time.Minute)
	return users, int32(total), nil
}

func (us *usersService) GetUserByUUID(ctx *gin.Context, user_uuid uuid.UUID) (sqlc.User, error) {
	// uid := uuid.MustParse(user_uuid)
	context := ctx.Request.Context()
	res, err := us.userRepo.GetUserByUUID(context, user_uuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return res, nil
}

func (us *usersService) AddUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {

	context := ctx.Request.Context()
	input.UserEmail = utils.NormalizeString(input.UserEmail)
	hashPass, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, &utils.AppError{Code: utils.ErrorCodeInternal, Message: "Hash password error"}
	}
	input.UserPassword = string(hashPass)
	usr, e := us.userRepo.AddUser(context, input)
	if e != nil {
		var pgEr *pgconn.PgError
		if errors.As(e, &pgEr) && (pgEr.Code == "23505") {
			return sqlc.User{}, utils.NewError("User already exists", utils.ErrorCodeConflict) //{Code: utils.ErrorCodeConflict, Message: "User email already exists"}
		}
		return sqlc.User{}, &utils.AppError{Code: utils.ErrorCodeInternal, Message: e.Error()}
	}
	//clear cache
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return usr, nil
}

func (us *usersService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserByUUIDParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.UserPassword != nil && *input.UserPassword != "" {
		hashPass, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, &utils.AppError{Code: utils.ErrorCodeInternal, Message: "Hash password error"}
		}
		hash := string(hashPass)
		input.UserPassword = &hash

	}
	res, err := us.userRepo.UpdateUser(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrorCodeNotFound)
		}
		return sqlc.User{}, &utils.AppError{Code: utils.ErrorCodeInternal, Message: "Update user error"}
	}

	//clear cache
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return res, nil
}
func (us *usersService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	rs, err := us.userRepo.SoftDeleteUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not found", utils.ErrorCodeNotFound)
		}
		return sqlc.User{}, utils.NewWrapError("Failed to delete user", utils.ErrorCodeInternal, err)
	}
	//clear cache
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return rs, nil
}
func (us *usersService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	rs, err := us.userRepo.Restore(context, uuid)
	if err != nil {

		return sqlc.User{}, utils.NewWrapError("Failed to restore user was not marked as deleted", utils.ErrorCodeInternal, err)
	}
	//clear cache
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return rs, nil
}
func (us *usersService) DeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()
	_, err := us.userRepo.Restore(context, uuid)
	if err != nil {

		return utils.NewWrapError("Failed to remove user", utils.ErrorCodeInternal, err)
	}
	//clear cache
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return nil

}

func (us *usersService) generateCacheKey(search string, order_by, sort string, page, limit int32, deleted bool) string {
	search = strings.TrimSpace(search)
	if search == "" {
		search = "none"
	}
	order_by = strings.TrimSpace(order_by)
	if order_by == "" {
		order_by = "user_id"
	}
	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}

	return fmt.Sprintf("users:%s:%s:%s:%d:%d:%t", search, order_by, sort, page, limit, deleted)
}
