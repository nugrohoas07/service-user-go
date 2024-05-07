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

// TODO
// encrypt password
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

func (usecase *usersUseCase) UpdateUserById(paramUserId string, updatedUser usersDto.UpdateUserRequest) error {
	// check is user id from param and payload match
	if paramUserId != updatedUser.ID {
		return fmt.Errorf("id not match")
	}
	// check is user id valid
	oldUserData, err := usecase.usersRepo.GetUserById(updatedUser.ID)
	if err != nil {
		return err
	}
	// edit process
	err = usecase.usersRepo.EditUser(oldUserData, updatedUser)
	if err != nil {
		return err
	}
	return nil
}

// TODO
// error if user delete itself
func (usecase *usersUseCase) DeleteUserById(userId string) error {
	// check is user id valid
	_, err := usecase.usersRepo.GetUserById(userId)
	if err != nil {
		return err
	}
	// soft deleting
	err = usecase.usersRepo.SoftDeleteUser(userId)
	if err != nil {
		return err
	}
	return nil
}
