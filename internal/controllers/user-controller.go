package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
	"app/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) CreateUser(ctx echo.Context) error {
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.CreateUserRequest)
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	// TODO: Perform additional validation on USN/Application Number and phone number formats
	success, err := uc.userService.CreateUser(ctx.Request().Context(), userID, reqBody)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}

	if success == false {
		return ctx.JSON(http.StatusConflict, map[string]string{
			"error": "user already exists",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}

func (uc *UserController) GetUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	user, err := uc.userService.GetUserProfile(ctx.Request().Context(), userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch user profile",
		})
	}

	if user == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.UpdateUserProfileRequest)

	// TODO: Perform additional validation on USN/Application Number and phone number formats

	success, err := uc.userService.UpdateUserProfile(ctx.Request().Context(), userID, reqBody)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update user profile",
		})
	}

	if !success {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "user profile update failed",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}
