package routes

import (
	"user-management-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	handler *handlers.UsersHandler
}

func NewUserRoute(h *handlers.UsersHandler) *UserRoute {
	return &UserRoute{handler: h}
}
func (ur *UserRoute) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/", ur.handler.GetUsers)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.POST("/", ur.handler.AddUser)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.DeleteUser)
	}
}
