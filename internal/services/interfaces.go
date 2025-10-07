package services

type UsersService interface {
	GetUsers()
	GetUserByUUID()
	AddUser()
	UpdateUser()
	DeleteUser()
}
