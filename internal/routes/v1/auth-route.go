package v1

import (
	handlers_v1 "user-management-api/internal/handlers/v1"

	"github.com/gin-gonic/gin"
)

type AuthRoute struct {
	//services services.AuthService
	handler *handlers_v1.AuthHandler
}

func NewAuthRoute(h *handlers_v1.AuthHandler) *AuthRoute {
	return &AuthRoute{handler: h}
}
func (ur *AuthRoute) Register(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("", ur.handler.Login)
		auth.POST("/logout", ur.handler.Logout)
		auth.POST("/refresh-token", ur.handler.RefreshToken)
	}
}
