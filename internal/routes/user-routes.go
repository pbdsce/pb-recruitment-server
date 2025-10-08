package routes

import (
	"app/internal/controllers"
	"app/internal/middleware"
	"app/internal/models/dto"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

func AddUserRoutes(
	e *echo.Echo,
	authClient *auth.Client,
	userController *controllers.UserController,
) {
	e.PUT("/users/create",
		userController.CreateUser,
		middleware.RequireFirebaseAuth(authClient),
		middleware.ValidateRequest(new(dto.CreateUserRequest)),
	)

	e.GET("/users/profile",
		userController.GetUserProfile,
		middleware.RequireFirebaseAuth(authClient),
	)

	e.POST("/users/profile",
		userController.UpdateUserProfile,
		middleware.RequireFirebaseAuth(authClient),
		middleware.ValidateRequest(new(dto.UpdateUserProfileRequest)),
	)
}
