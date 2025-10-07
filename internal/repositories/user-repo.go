package repositories

import (
	"log"
	"user-management-api/internal/models"
)

type UserRepoInMemory struct {
	users []models.Users
}

func NewUserRepo() UserRepo {
	return &UserRepoInMemory{
		users: make([]models.Users, 0),
	}
}

func (ur *UserRepoInMemory) GetUsers() {
	log.Printf("GetUsers in UserRepoInMemory")
}

func (ur *UserRepoInMemory) AddUser() {}

func (ur *UserRepoInMemory) GetUserByUUID() {}

func (ur *UserRepoInMemory) UpdateUser() {}

func (ur *UserRepoInMemory) DeleteUser() {}
