package usersDelivery

import (
	"fmt"
	"service-user/middlewares"
	"service-user/model/dto/json"
	"service-user/model/dto/usersDto"
	"service-user/pkg/validation"
	"service-user/src/users"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

type usersDelivery struct {
	usersUC users.UsersUseCase
}

func NewUsersDelivery(v1Group *gin.RouterGroup, usersUC users.UsersUseCase) {
	handler := usersDelivery{
		usersUC: usersUC,
	}
	store := memstore.NewStore([]byte("secret"))
	usersGroup := v1Group.Group("/users")
	{
		usersGroup.POST("/login", middlewares.BasicAuth, sessions.Sessions("mysession", store), handler.Login)
		usersGroup.POST("/create", middlewares.JwtAuth(), handler.AddUser)                                          // create new user
		usersGroup.GET("/", middlewares.JwtAuth(), handler.GetUsers)                                                //get list all users
		usersGroup.GET("/:id", middlewares.JwtAuth(), handler.GetUserById)                                          // get user data by userId
		usersGroup.PUT("/:id", middlewares.JwtAuth(), handler.UpdateUser)                                           // edit user data by userId
		usersGroup.DELETE("/:id", middlewares.JwtAuth(), sessions.Sessions("mysession", store), handler.DeleteUser) // soft delete user by userId
	}
}

func (ud *usersDelivery) Login(ctx *gin.Context) {
	var req usersDto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	err := ud.usersUC.ValidateEmailPass(req.Email, req.Password)
	if err != nil {
		json.NewAbortUnauthorized(ctx, err.Error(), "01", "01")
		return
	}

	token, err := middlewares.GenerateTokenJwt(req.Email, 3)
	if err != nil {
		errMsg := "internal server error"
		json.NewResponseError(ctx, errMsg, "01", "01")
		return
	}

	userId, err := ud.usersUC.GetUserIdByEmail(req.Email)
	if err != nil {
		json.NewAbortUnauthorized(ctx, err.Error(), "01", "01")
		return
	}

	session := sessions.Default(ctx)
	session.Set("userId", userId)
	session.Save()

	json.NewResponseSuccess(ctx, map[string]interface{}{"token": token}, "", "01", "01")
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

func (ud *usersDelivery) GetUsers(ctx *gin.Context) {
	var queryParams usersDto.Query
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	listUsers, totalData, err := ud.usersUC.GetAllUsers(queryParams)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}

	if len(listUsers) == 0 {
		json.NewResponseSuccess(ctx, nil, "data not found", "01", "01")
		return
	}

	var paging json.Paging
	if queryParams.Page != 0 && queryParams.Size != 0 {
		paging = json.Paging{
			Page:      queryParams.Page,
			TotalData: totalData,
		}
	}

	json.NewResponseSuccessWithPaging(ctx, listUsers, paging, "", "01", "02")
}

func (ud *usersDelivery) GetUserById(ctx *gin.Context) {
	var param usersDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}
	userData, err := ud.usersUC.GetUserById(param.ID)
	if err != nil {
		json.NewResponseNotFound(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponseSuccess(ctx, userData, "", "01", "01")
}

func (ud *usersDelivery) UpdateUser(ctx *gin.Context) {
	var param usersDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	var userUpdatePayload usersDto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&userUpdatePayload); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	err := ud.usersUC.UpdateUserById(param.ID, userUpdatePayload)
	if err != nil {
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}

	json.NewResponseSuccess(ctx, nil, "success", "01", "02")
}

func (ud *usersDelivery) DeleteUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	loggedUserId := session.Get("userId")
	fmt.Println("DARI SESSION :", loggedUserId)

	var param usersDto.Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		validationError := validation.GetValidationError(err)
		if len(validationError) > 0 {
			json.NewResponseBadRequest(ctx, validationError, "bad request", "01", "02")
			return
		}
	}

	if loggedUserId == param.ID {
		json.NewAbortForbidden(ctx, "action not allowed", "01", "01")
		return
	}

	err := ud.usersUC.DeleteUserById(param.ID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			json.NewResponseNotFound(ctx, err.Error(), "01", "01")
			return
		}
		json.NewResponseError(ctx, err.Error(), "01", "01")
		return
	}
	json.NewResponseSuccess(ctx, nil, "success", "01", "02")
}
