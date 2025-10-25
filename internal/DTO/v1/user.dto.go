package dto_v1

import (
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"

	"github.com/google/uuid"
)

type UserDto struct {
	UUID      string `json:"uuid" `
	Name      string `json:"name" `
	Email     string `json:"email" `
	Age       *int   `json:"age" ` // maybe is nil
	Status    string `json:"status" `
	Level     string `json:"level"`
	CreatedAt string `json:"created_at"`
	// UpdatedAt string `json:"updated_at"`
	// DeletedAt string `json:"deleted_at" `
}
type GetUserByUUIDParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}
type GetUsersByQuery struct {
	Search string `form:"search" binding:"omitempty,search_format"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
	Order  string `form:"order_by" binding:"omitempty,oneof=user_id user_created_at"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

type CreateUsersInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Age      int32  `json:"age" binding:omitempty,gt=18"`
	Password string `json:"password" binding:"required,min=6,password_strong"`
	Status   int32  `json:"status" binding:"required,oneof=1 2 3"`
	Level    int32  `json:"level" binding:"required,oneof=1 2 3"`
}

type UpdateUsersInput struct {
	Name *string `json:"name" binding:"omitempty"`
	// Email    string  `json:"email" binding:"required,email,email_advanced"`
	Age      *int32  `json:"age" binding:"omitempty,gt=18"`
	Password *string `json:"password" binding:"omitempty,min=6,password_strong"`
	Status   *int32  `json:"status" binding:"omitempty,oneof=1 2 3"`
	Level    *int32  `json:"level" binding:"omitempty,oneof=1 2 3"`
}

func (input *CreateUsersInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		UserEmail:    input.Email,
		UserFullname: input.Name,
		UserAge:      utils.ConvertInt32ToPointer(input.Age), //input.Age,
		UserPassword: input.Password,
		UserStatus:   int32(input.Status),
		UserLevel:    int32(input.Level),
	}
}

func (input *UpdateUsersInput) MapUpdateInputToModel(userUuid uuid.UUID) sqlc.UpdateUserByUUIDParams {
	return sqlc.UpdateUserByUUIDParams{
		UserUuid:     userUuid,
		UserFullname: input.Name,
		UserAge:      (input.Age), //input.Age,
		UserPassword: input.Password,
		UserStatus:   (input.Status),
		UserLevel:    (input.Level),
	}
}

func MapUserToDto(u sqlc.User) *UserDto {
	dtos := &UserDto{
		UUID:  u.UserUuid.String(),
		Name:  u.UserFullname,
		Email: u.UserEmail,
		// Age:    int(*u.UserAge),
		Status:    mapStatusText(int(u.UserStatus)),
		Level:     mapLevelText(int(u.UserLevel)),
		CreatedAt: u.UserCreatedAt.Format("2006-01-01 15:04:01"),
	}
	if u.UserAge != nil { //check nil
		age := int(*u.UserAge)
		dtos.Age = &age
	}
	// if u.UserDeletedAt.Valid { //check nil
	// 	dtos.DeletedAt = u.UserDeletedAt.Time.Format("2006-01-01 15:04:01")
	// } else {
	// 	dtos.DeletedAt = ""
	// }
	return dtos
}

//	func MapUsersToDto() *UserDto {
//		return &UserDto{}
//	}
func MapUsersToDto(users []sqlc.User) []UserDto {
	dtos := make([]UserDto, 0, len(users))
	for _, u := range users {
		dtos = append(dtos, *MapUserToDto(u))
	}
	return dtos
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "ACTIVE"
	case 2:
		return "INACTIVE"
	case 3:
		return "BANNED"
	default:
		return "None"
	}
}

func mapLevelText(level int) string {
	switch level {
	case 1:
		return "ADMINTRATOR"
	case 2:
		return "MORDERATOR"
	case 3:
		return "MEMBER"
	default:
		return "None"
	}
}
