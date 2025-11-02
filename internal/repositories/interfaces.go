<<<<<<< HEAD
package repositories

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUsersV2(ctx context.Context, search *string, orderBy, sort string, offset, limit int32) ([]sqlc.User, error)
	GetUsers(ctx context.Context, search *string, orderBy, sort string, offset, limit int32) ([]sqlc.User, error)
	UsersCount(ctx context.Context, search *string) (int64, error)
	AddUser(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	UpdateUser(ctx context.Context, input sqlc.UpdateUserByUUIDParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	FindByEmail(email string)
}
=======
package repositories

import (
	"context"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepo interface {
	GetUsersV2(ctx context.Context, search *string, orderBy, sort string, offset, limit int32, deleted bool) ([]sqlc.User, error)
	GetUsers(ctx context.Context, search *string, orderBy, sort string, offset, limit int32) ([]sqlc.User, error)
	UsersCount(ctx context.Context, search *string, deleted bool) (int64, error)
	AddUser(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	UpdateUser(ctx context.Context, input sqlc.UpdateUserByUUIDParams) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	FindByEmail(email string)
}
>>>>>>> 1bd3d85b166d78e8ef8b54770c445ebfac40b114
