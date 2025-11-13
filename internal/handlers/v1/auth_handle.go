package handlers_v1

import (
	"net/http"
	dto_v1 "user-management-api/internal/DTO/v1"
	"user-management-api/internal/services"
	"user-management-api/internal/utils"
	"user-management-api/internal/validations"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	//services services.AuthService
	service services.AuthService
}

func NewAuthHandler(s services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: s,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	// Implementation for login

	var input dto_v1.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	accessToken, refreshToken, expiresIn, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	response := dto_v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Login successful", response)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	// Implementation for logout
	var input dto_v1.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	err := ah.service.Logout(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Logout successfully")
}

func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	// Implementation for logout
	var input dto_v1.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	accessToken, refreshToken, expiresIn, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	response := dto_v1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Refresh token generate successfully", response)
}

func (ah *AuthHandler) RequestForgotPassword(ctx *gin.Context) {
	var input dto_v1.RequestPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	if err := ah.service.RequestForgotPassword(ctx, input.Email); err != nil {
		utils.NewError("", utils.ErrorCodeInternal)
		return
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Reset link is sent to email !")
}
func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var input dto_v1.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	if err := ah.service.ResetPassword(ctx, input.Token, input.NewPassword); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Password reset successfully")
}
