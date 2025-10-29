package handlers_v1

import (
	"log"
	"net/http"
	dto_v1 "user-management-api/internal/DTO/v1"
	"user-management-api/internal/services"
	"user-management-api/internal/utils"
	"user-management-api/internal/validations"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsersHandler struct {
	services services.UsersService
}

type GetUsersByQuery struct {
	Search string `form:"search" binding:"omitempty,search_format"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
}

func NewUsersHandler(sv services.UsersService) *UsersHandler {

	return &UsersHandler{
		services: sv,
	}
}

func (uh *UsersHandler) GetUsers(ctx *gin.Context) {
	var query dto_v1.GetUsersByQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	log.Printf("query: %+v", query)
	//users, total, err := uh.services.GetUsers(ctx, &query.Search, query.Order, query.Sort, query.Page, query.Limit, query.Deleted)
	users, total, err := uh.services.GetUsers(ctx, &query.Search, query.Order, query.Sort, query.Page, query.Limit, false)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	////
	dtos := dto_v1.MapUsersToDto(users) // convert model to list dto for response
	paginationResponse := utils.NewPaginationResponse(dtos, query.Page, query.Limit, total)
	utils.ResponseSuccess(ctx, http.StatusOK, "Users retrieved successfully", paginationResponse)
}

func (uh *UsersHandler) GetUserByUUID(ctx *gin.Context) {
	var params dto_v1.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	//log.Printf("uuid: %+v", params.Uuid)
	uid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	usr, er := uh.services.GetUserByUUID(ctx, uid)
	if er != nil {
		utils.ResponseError(ctx, er)
		return
	}
	dto := dto_v1.MapUserToDto(usr)
	utils.ResponseSuccess(ctx, http.StatusOK, "User retrieved successfully", dto)
}
func (uh *UsersHandler) GetUsersDeleted(ctx *gin.Context) {

	var query dto_v1.GetUsersByQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	log.Printf("query: %+v", query)

	users, total, err := uh.services.GetUsers(ctx, &query.Search, query.Order, query.Sort, query.Page, query.Limit, true)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	////
	dtos := dto_v1.MapUsersToDto(users) // convert model to list dto for response
	paginationResponse := utils.NewPaginationResponse(dtos, query.Page, query.Limit, total)
	utils.ResponseSuccess(ctx, http.StatusOK, "Users soft deleted successfully", paginationResponse)
}
func (uh *UsersHandler) AddUser(ctx *gin.Context) {
	var input dto_v1.CreateUsersInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	u := input.MapCreateInputToModel()
	usr, err := uh.services.AddUser(ctx, u)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	dto := dto_v1.MapUserToDto(usr)
	utils.ResponseSuccess(ctx, http.StatusCreated, "User created successfully", dto)
	//ctx.JSON(http.StatusOK, )
}

func (uh *UsersHandler) UpdateUser(ctx *gin.Context) {
	var params dto_v1.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	// log.Printf("uuid: %+v", params.Uuid)
	var input dto_v1.UpdateUsersInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	model := input.MapUpdateInputToModel(userUuid)

	usr, er := uh.services.UpdateUser(ctx, model)
	if er != nil {
		utils.ResponseError(ctx, er)
		return
	}
	response := dto_v1.MapUserToDto(usr) // convert model to dto
	utils.ResponseSuccess(ctx, http.StatusOK, "User updated successfully", response)
}
func (uh *UsersHandler) SoftDeleteUser(ctx *gin.Context) {

	var params dto_v1.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	usr, err := uh.services.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	response := dto_v1.MapUserToDto(usr) // convert model to dto
	utils.ResponseSuccess(ctx, http.StatusOK, "User deleted successfully", response)
}
func (uh *UsersHandler) DeleteUser(ctx *gin.Context) {

	var params dto_v1.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	e := uh.services.DeleteUser(ctx, userUuid)
	if e != nil {
		utils.ResponseError(ctx, e)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
func (uh *UsersHandler) RestoreUser(ctx *gin.Context) {

	var params dto_v1.GetUserByUUIDParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validations.HandleValidationErr(err))
		return
	}
	userUuid, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	usr, err := uh.services.RestoreUser(ctx, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	response := dto_v1.MapUserToDto(usr) // convert model to dto

	utils.ResponseSuccess(ctx, http.StatusOK, "User restored successfully", response)
}
func (uh *UsersHandler) GetUsersByRole(ctx *gin.Context) {

}
