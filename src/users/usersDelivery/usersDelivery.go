package usersDelivery

import (
	"service-user/model/dto/json"
	"service-user/model/dto/usersDto"
	"service-user/pkg/validation"
	"service-user/src/users"

	"github.com/gin-gonic/gin"
)

type usersDelivery struct {
	usersUC users.UsersUseCase
}

func NewUsersDelivery(v1Group *gin.RouterGroup, usersUC users.UsersUseCase) {
	handler := usersDelivery{
		usersUC: usersUC,
	}
	usersGroup := v1Group.Group("/users")
	{
		usersGroup.POST("/login")                   // login user with email:password
		usersGroup.POST("/create", handler.AddUser) // create new user
		usersGroup.GET("/")                         //get list all users
		usersGroup.GET("/:id")                      // get user data by userId
		usersGroup.PUT("/:id")                      // edit user data by userId
		usersGroup.DELETE("/:id")                   // soft delete user by userId
	}
}

func (ud *usersDelivery) AddUser(ctx *gin.Context) {
	var userPayload usersDto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userPayload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	err := ud.usersUC.AddUser(userPayload)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success", "01", "01")
}
