package services

import (
	"log"
	"user-management-api/internal/repositories"
)

type usersService struct {
	userRepo repositories.UserRepo
}

func NewUsersService(repo repositories.UserRepo) UsersService {
	return &usersService{
		userRepo: repo,
	}
}

func (us *usersService) GetUsers() {
	log.Printf("GetUsers in usersService")
	us.userRepo.GetUsers()
}

func (us *usersService) GetUserByUUID() {}

func (us *usersService) AddUser() {}

func (us *usersService) UpdateUser() {}

func (us *usersService) DeleteUser() {}
