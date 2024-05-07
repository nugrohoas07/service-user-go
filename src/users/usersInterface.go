package users

import "service-user/model/dto/usersDto"

type UsersRepository interface {
	InsertUser(newUser usersDto.CreateUserRequest) error
	GetUserPassword(email string) (string, error)
	CheckEmailExist(email string) string
}

type UsersUseCase interface {
	Login(email, password string) error
	AddUser(newUser usersDto.CreateUserRequest) error
}
