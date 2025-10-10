package controllers

import (
	"app/internal/common"
	"app/internal/models/dto"
	"app/internal/services"
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

// validation
func (uc *UserController) UpdateUserProfile(ctx echo.Context) error {
	userID := ctx.Get(common.AUTH_USER_ID).(string)
	reqBody := ctx.Get(common.VALIDATED_REQUEST_BODY).(*dto.UpdateUserProfileRequest)

	usnPattern := `^[0-9][A-Z]{2}[0-9]{2}[A-Z]{2}[0-9]{3}$`  // Example: 1DS24CS123
	appNumberPattern := `^[0-9]{2}[A-Z]{5}[0-9]{4}$`         // Example: 24UGDSI1650
	phonePattern := `^\+91[0-9]{10}$`                        // Example: +911234567890

	usnUpper := strings.ToUpper(reqBody.USN)
	isUSN, _ := regexp.MatchString(usnPattern, usnUpper)
	isAppNumber, _ := regexp.MatchString(appNumberPattern, usnUpper)

	// Validation of USN or Application Number based on CurrentYear
	if reqBody.CurrentYear == 1 {
		if !isAppNumber {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "First-year students must provide a valid Application Number (e.g., 24UGDSI1650)",
			})
		}
	} else {
		if !isUSN {
			return ctx.JSON(http.StatusBadRequest, map[string]string{
				"error": "Provide a valid USN (e.g., 1DS24CS123)",
			})
		}
	}

	// mob no. validation
	if matched, _ := regexp.MatchString(phonePattern, reqBody.MobileNumber); !matched {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid mobile number format. Expected format: +91XXXXXXXXXX",
		})
	}

	if err := uc.userService.UpdateUserProfile(userID, reqBody); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to update user profile",
		})
	}

	return ctx.NoContent(http.StatusCreated)
}
