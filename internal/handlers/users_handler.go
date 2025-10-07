package handlers

import (
	"log"
	"user-management-api/internal/services"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	services services.UsersService
}

func NewUsersHandler(sv services.UsersService) *UsersHandler {
	return &UsersHandler{
		services: sv,
	}
}

func (uh *UsersHandler) GetUsers(ctx *gin.Context) {
	log.Printf("GetUsers in UsersHandler")
	uh.services.GetUsers()
}

func (uh *UsersHandler) GetUserByUUID(ctx *gin.Context) {}

func (uh *UsersHandler) AddUser(ctx *gin.Context) {}

func (uh *UsersHandler) UpdateUser(ctx *gin.Context) {}

func (uh *UsersHandler) DeleteUser(ctx *gin.Context) {}
