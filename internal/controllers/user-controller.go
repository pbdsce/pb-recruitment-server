package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
	"app/internal/services"
	"net/http"
	"fmt"
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
	usnPattern := `^[0-9][A-Z]{2}[0-9]{2}[A-Z]{2}[0-9]{3}$`  // Example: 1DS24CS123
	appNumberPattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}$`         // Example: 24UGDSI1650
	phonePattern := `^\+91[0-9]{10}$`                        // Example: +911234567890

	usnUpper := strings.ToUpper(usn)
	isUSN, _ := regexp.MatchString(usnPattern, usnUpper)
	isAppNumber, _ := regexp.MatchString(appNumberPattern, usnUpper)

	if currentYear == 1 && !isAppNumber {
		return fmt.Errorf("First-year students must provide a valid Application Number (e.g., 24UGDSI1650)")
	} else if currentYear != 1 && !isUSN {
		return fmt.Errorf("Provide a valid USN (e.g., 1DS24CS123)")
	}
	if matched, _ := regexp.MatchString(phonePattern, mobile); !matched {
		return fmt.Errorf("Invalid mobile number format. Expected format: +91XXXXXXXXXX")
	}
	return nil
}

func (uc *UserController) CreateUser(ctx echo.Context) error {
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.CreateUserRequest)
	userID := ctx.Get(common.AUTH_USER_ID).(string)

	if err := validateUserInput(reqBody.USN, reqBody.MobileNumber, reqBody.CurrentYear); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

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

	if err := validateUserInput(reqBody.USN, reqBody.MobileNumber, reqBody.CurrentYear); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

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
