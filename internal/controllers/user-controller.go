package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
	"app/internal/services"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

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

func validateUserInput(usn string, mobile string, currentYear int) error {
	usnPattern := `^1DS2[3-4](AI|AE|AU|BT|CG|MD|ET|EC|ME|EE|CH|CB|IC|CD|CS|CV|IS)[0-9]{3}$` // Example: 1DS24IC015
	appNumberPattern := `^25UGDS[0-9]{4}$`                                                  // Example: 25UGDS1234
	phonePattern := `^[0-9]{10}$`                                                           // Example: 9234567890

	usnUpper := strings.ToUpper(usn)

	if currentYear == 1 {
		if matched, _ := regexp.MatchString(appNumberPattern, usnUpper); !matched {
			return fmt.Errorf("first-year students must provide a valid Application Number (e.g., 25UGDS1234)")
		}
	} else {
		if matched, _ := regexp.MatchString(usnPattern, usnUpper); !matched {
			return fmt.Errorf("provide a valid USN (e.g., 1DS24IC015)")
		}
	}

	if matched, _ := regexp.MatchString(phonePattern, mobile); !matched {
		return fmt.Errorf("invalid mobile number format")
	}
	return nil
}

func (uc *UserController) CreateUser(ctx echo.Context) error {
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.CreateUserRequest)
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	if err := validateUserInput(reqBody.USN, reqBody.MobileNumber, reqBody.CurrentYear); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := uc.userService.CreateUser(ctx.Request().Context(), userID, reqBody); err != nil {
		if errors.Is(err, common.UserAlreadyExistsError{}) {
			return ctx.JSON(http.StatusConflict, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}

func (uc *UserController) GetUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	user, err := uc.userService.GetUserProfile(ctx.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, common.UserNotFoundError{}) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch user profile",
		})
	}

	return ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.UpdateUserProfileRequest)

	// mob no. validation
	phonePattern := `^[0-9]{10}$`
	if matched, _ := regexp.MatchString(phonePattern, reqBody.MobileNumber); !matched {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid mobile number format",
		})
	}

	if err := uc.userService.UpdateUserProfile(ctx.Request().Context(), userID, reqBody); err != nil {
		if errors.Is(err, common.UserNotFoundError{}) {
			return ctx.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update user profile",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}
