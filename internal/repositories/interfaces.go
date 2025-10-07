package repositories

type UserRepo interface {
	GetUsers()
	AddUser()
	GetUserByUUID()
	UpdateUser()
	DeleteUser()
}
