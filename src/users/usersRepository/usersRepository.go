package usersRepository

import (
	"database/sql"
	"fmt"
	"service-user/model/dto/usersDto"
	"service-user/model/entity/usersEntity"
	"service-user/src/users"
)

type usersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) users.UsersRepository {
	return &usersRepository{db}
}

func (repo *usersRepository) InsertUser(newUser usersDto.CreateUserRequest) error {
	query := "INSERT INTO users (fullname,email,password) VALUES($1, $2, $3)"
	_, err := repo.db.Exec(query, newUser.FullName, newUser.Email, newUser.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *usersRepository) GetUserById(userId string) (usersEntity.UserData, error) {
	var userData usersEntity.UserData
	query := "SELECT id,fullname,email,password FROM users WHERE id = $1"
	err := repo.db.QueryRow(query, userId).Scan(&userData.ID, &userData.FullName, &userData.Email, &userData.Password)
	if err != nil {
		return usersEntity.UserData{}, err
	}
	return userData, nil
}

func (repo *usersRepository) GetUserPassword(email string) (string, error) {
	var storedPassword string

	query := "SELECT password FROM users WHERE email = $1"
	err := repo.db.QueryRow(query, email).Scan(&storedPassword)
	if err != nil {
		return "", err
	}

	return storedPassword, nil
}

func (repo *usersRepository) CheckEmailExist(email string) string {
	var duplicateEmail string
	query := "SELECT email FROM users WHERE email = $1"
	repo.db.QueryRow(query, email).Scan(&duplicateEmail)
	fmt.Println("email didapat :", duplicateEmail)
	return duplicateEmail
}
