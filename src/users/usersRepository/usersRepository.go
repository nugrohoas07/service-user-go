package usersRepository

import (
	"database/sql"
	"fmt"
	"service-user/model/dto/usersDto"
	"service-user/model/entity/usersEntity"
	"service-user/src/users"
	"strings"
	"time"
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

func (repo *usersRepository) GetUsers(queryParams usersDto.Query) ([]usersEntity.UserData, int, error) {
	fmt.Println("MASUK REPO")
	var rows *sql.Rows
	var err error

	query := `SELECT id,fullname,email,password FROM users
	WHERE 1=1
	AND ($1 = '' OR email = $1)
	AND ($2 = '' OR fullname = $2)
	ORDER BY fullname ASC`

	countQuery := `SELECT COUNT(*)
	FROM users
	WHERE 1=1
	AND ($1 = '' OR email = $1)
	AND ($2 = '' OR fullname = $2)`

	if queryParams.Page != 0 && queryParams.Size != 0 {
		offset := (queryParams.Page - 1) * queryParams.Size
		query += " LIMIT $3 OFFSET $4"
		rows, err = repo.db.Query(query, queryParams.Email, queryParams.Fullname, queryParams.Size, offset)
	} else {
		rows, err = repo.db.Query(query, queryParams.Email, queryParams.Fullname)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	//get total data
	var totalData int
	err = repo.db.QueryRow(countQuery, queryParams.Email, queryParams.Fullname).Scan(&totalData)
	if err != nil {
		fmt.Println("ERROR PAS NGITUNG")
		return nil, 0, fmt.Errorf("internal server error")
	}

	listUsers := scanUsers(rows)
	fmt.Println("list :", listUsers)
	fmt.Println("totalData :", totalData)
	return listUsers, totalData, nil
}

func (repo *usersRepository) GetUserById(userId string) (usersEntity.UserData, error) {
	var userData usersEntity.UserData
	query := "SELECT id,fullname,email,password FROM users WHERE id = $1 AND deleted_at IS NULL"
	err := repo.db.QueryRow(query, userId).Scan(&userData.ID, &userData.FullName, &userData.Email, &userData.Password)
	if err != nil {
		return usersEntity.UserData{}, fmt.Errorf("user not found")
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

func (repo *usersRepository) EditUser(oldUser usersEntity.UserData, updatedUser usersDto.UpdateUserRequest) error {
	if strings.TrimSpace(updatedUser.FullName) == "" {
		updatedUser.FullName = oldUser.FullName
	}
	if strings.TrimSpace(updatedUser.Password) == "" {
		updatedUser.Password = oldUser.Password
	}
	query := "UPDATE users SET fullname = $1, password = $2 WHERE id = $3"
	_, err := repo.db.Exec(query, updatedUser.FullName, updatedUser.Password, updatedUser.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *usersRepository) SoftDeleteUser(userId string) error {
	query := "UPDATE users SET deleted_at = $1 WHERE id = $2"
	_, err := repo.db.Exec(query, time.Now(), userId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *usersRepository) CheckEmailExist(email string) string {
	var duplicateEmail string
	query := "SELECT email FROM users WHERE email = $1"
	repo.db.QueryRow(query, email).Scan(&duplicateEmail)
	return duplicateEmail
}

func scanUsers(rows *sql.Rows) []usersEntity.UserData {
	var users []usersEntity.UserData
	var err error
	for rows.Next() {
		user := usersEntity.UserData{}
		err = rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Password)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	return users
}
