package router

import (
	"database/sql"
	"service-user/src/users/usersDelivery"
	"service-user/src/users/usersRepository"
	"service-user/src/users/usersUseCase"

	"github.com/gin-gonic/gin"
)

func InitRoute(v1Group *gin.RouterGroup, db *sql.DB) {
	usersRepo := usersRepository.NewUsersRepository(db)
	usersUseCase := usersUseCase.NewUsersUseCase(usersRepo)
	usersDelivery.NewUsersDelivery(v1Group, usersUseCase)
}
