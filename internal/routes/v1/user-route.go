package v1

import (
	handlers_v1 "user-management-api/internal/handlers/v1"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	handler *handlers_v1.UsersHandler
}

func NewUserRoute(h *handlers_v1.UsersHandler) *UserRoute {
	return &UserRoute{handler: h}
}
func (ur *UserRoute) Register(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.GET("/", ur.handler.GetUsers)
		users.GET("/:uuid", ur.handler.GetUserByUUID)
		users.GET("/soft_deleted", ur.handler.GetUsersDeleted)
		users.POST("/", ur.handler.AddUser)
		users.PUT("/:uuid", ur.handler.UpdateUser)
		users.DELETE("/:uuid", ur.handler.SoftDeleteUser)
		users.PUT("/restore/:uuid", ur.handler.RestoreUser)
		users.DELETE("/trash/:uuid", ur.handler.DeleteUser)
	}
}
