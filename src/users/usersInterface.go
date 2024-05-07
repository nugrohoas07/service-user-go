package users

import (
	"service-user/model/dto/usersDto"
	"service-user/model/entity/usersEntity"
)

type UsersRepository interface {
	InsertUser(newUser usersDto.CreateUserRequest) error
	GetUserPassword(email string) (string, error)
	CheckEmailExist(email string) string
	GetUserById(userId string) (usersEntity.UserData, error)
	SoftDeleteUser(userId string) error
	EditUser(oldUser usersEntity.UserData, updatedUser usersDto.UpdateUserRequest) error
}

type UsersUseCase interface {
	Login(email, password string) error
	AddUser(newUser usersDto.CreateUserRequest) error
	GetUserById(userId string) (usersEntity.UserData, error)
	DeleteUserById(userId string) error
	UpdateUserById(paramUserId string, updatedUser usersDto.UpdateUserRequest) error
}
