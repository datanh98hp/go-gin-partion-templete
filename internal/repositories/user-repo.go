package repositories

import (
	"context"
	"fmt"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository struct {
	db sqlc.Querier
}

func NewUserRepo(db sqlc.Querier) UserRepo {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUsersV2(ctx context.Context, search *string, orderBy, sort string, offset, limit int32) ([]sqlc.User, error) {
	query := `SELECT * 
FROM users 
WHERE user_deleted_at IS NULL 
AND (
   	$1::TEXT IS NULL 
    OR $1::TEXT = ''
    OR user_email ILIKE '%' || $1 || '%' 
    OR user_fullname ILIKE '%'|| $1 || '%'
)`
	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}
	switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	default:
		query += " ORDER BY user_id ASC"
	}

	query += " LIMIT $2 OFFSET $3"

	rows, err := db.DBpool.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserFullname,
			&i.UserPassword,
			&i.UserAge,
			&i.UserStatus,
			&i.UserLevel,
			&i.UserDeletedAt,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

func (ur *UserRepository) GetUsers(ctx context.Context, search *string, orderBy, sort string, offset, limit int32) ([]sqlc.User, error) {

	var (
		users []sqlc.User
		er    error
	)

	switch {
	case orderBy == "user_id" && sort == "asc":
		users, er = ur.db.GetUsersIdAsc(ctx, sqlc.GetUsersIdAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_id" && sort == "desc":
		users, er = ur.db.GetUsersIdDesc(ctx, sqlc.GetUsersIdDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_created_at" && sort == "desc":
		users, er = ur.db.GetUsersCreatedAtDesc(ctx, sqlc.GetUsersCreatedAtDescParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})

	case orderBy == "user_created_at" && sort == "asc":
		users, er = ur.db.GetUsersCreatedAtAsc(ctx, sqlc.GetUsersCreatedAtAscParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})

		////
		if er != nil {
			return []sqlc.User{}, er
		}
		return users, nil
	}

	return users, nil
}

func (ur *UserRepository) AddUser(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	// log.Printf("userParam: %+v", input)
	u, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}
	return u, nil
}

func (ur *UserRepository) GetUserByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {

	usr, er := ur.db.GetUserByUUID(ctx, userUuid)
	if er != nil {
		return sqlc.User{}, er
	}
	return usr, nil

}

func (ur *UserRepository) UpdateUser(ctx context.Context, input sqlc.UpdateUserByUUIDParams) (sqlc.User, error) {
	urs, er := ur.db.UpdateUserByUUID(ctx, input)
	if er != nil {
		return sqlc.User{}, er
	}
	return urs, nil
}
func (ur *UserRepository) SoftDeleteUser(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	sdel, er := ur.db.SoftDeleteUserByUUID(ctx, uuid)
	if er != nil {
		return sqlc.User{}, er
	}
	return sdel, nil
}
func (ur *UserRepository) Restore(ctx context.Context, uuid uuid.UUID) (sqlc.User, error) {
	restore, er := ur.db.RestoreUsers(ctx, uuid)
	if er != nil {
		return sqlc.User{}, er
	}
	return restore, nil
}

func (ur *UserRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, er := ur.db.TrashUsers(ctx, uuid)
	if er != nil {
		return er
	}
	return nil
}

func (ur *UserRepository) FindByEmail(email string) {
}

func (ur *UserRepository) UsersCount(ctx context.Context, search *string) (int64, error) {
	total, er := ur.db.CountUsers(ctx, search)
	if er != nil {
		return 0, er
	}
	return total, nil

}
