package usersUseCase

import (
	"fmt"
	"service-user/model/dto/usersDto"
	"service-user/model/entity/usersEntity"
	"service-user/src/users"
)

type usersUseCase struct {
	usersRepo users.UsersRepository
}

func NewUsersUseCase(usersRepo users.UsersRepository) users.UsersUseCase {
	return &usersUseCase{usersRepo}
}

func (usecase *usersUseCase) Login(username, password string) error {
	storedPassword, err := usecase.usersRepo.GetUserPassword(username)
	fmt.Println("password di db :", storedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (usecase *usersUseCase) AddUser(newUser usersDto.CreateUserRequest) error {
	existedEmail := usecase.usersRepo.CheckEmailExist(newUser.Email)
	if existedEmail == newUser.Email {
		return fmt.Errorf("email already used")
	}
	err := usecase.usersRepo.InsertUser(newUser)
	if err != nil {
		return err
	}
	return nil
}

func (usercase *usersUseCase) GetUserById(userId string) (usersEntity.UserData, error) {
	userData, err := usercase.usersRepo.GetUserById(userId)
	if err != nil {
		return usersEntity.UserData{}, err
	}
	return userData, nil
}

// TODO
// error if user delete itself
func (usercase *usersUseCase) DeleteUserById(userId string) error {
	// check is user id valid
	_, err := usercase.usersRepo.GetUserById(userId)
	if err != nil {
		return err
	}
	err = usercase.usersRepo.SoftDeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}
