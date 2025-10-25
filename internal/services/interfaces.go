package services

import (
	"user-management-api/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersService interface {
	GetUsers(ctx *gin.Context, search *string, order_by, sort string, page, limit int32) ([]sqlc.User, int32, error)
	GetUserByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	AddUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, input sqlc.UpdateUserByUUIDParams) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error)
	DeleteUser(ctx *gin.Context, uuid uuid.UUID) error
}
