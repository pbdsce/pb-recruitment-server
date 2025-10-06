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

func (uc *UserController) GetUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	user, err := uc.userService.GetUserProfile(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch user profile",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.UpdateUserProfileRequest)

	// TODO: Perform additional validation on USN/Application Number and phone number formats

	if err := uc.userService.UpdateUserProfile(userID, reqBody); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update user profile",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}
