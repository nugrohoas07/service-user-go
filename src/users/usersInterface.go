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
	GetUsers(queryParams usersDto.Query) ([]usersEntity.UserData, int, error)
	GetUserIdByEmail(email string) (string, error)
}

type UsersUseCase interface {
	ValidateEmailPass(username, password string) error
	AddUser(newUser usersDto.CreateUserRequest) error
	GetUserById(userId string) (usersEntity.UserData, error)
	DeleteUserById(userId string) error
	UpdateUserById(paramUserId string, updatedUser usersDto.UpdateUserRequest) error
	GetAllUsers(queryParams usersDto.Query) ([]usersEntity.UserData, int, error)
	GetUserIdByEmail(email string) (string, error)
}
