package routes

import (
	"app/internal/controllers"
	"app/internal/middleware"
	"app/internal/services"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

func AddAdminRoutes(
	e *echo.Echo,
	contestController *controllers.ContestController,
	authClient *auth.Client,
	userService *services.UserService,
	adminService *services.AdminService,
) {
	adminGroup := e.Group("/admin")

	adminGroup.Use(middleware.RequireFirebaseAuth(authClient))

	adminGroup.Use(middleware.RequireAdminRole(userService, adminService))

	//Contest Management
	adminGroup.POST("/contest", contestController.HandleCreateContest)
	adminGroup.PUT("/contest/:id", contestController.HandleUpdateContest)
	adminGroup.DELETE("/contest/:id", contestController.HandleDeleteContest)

	//Problem Management
	adminGroup.POST("/:contestid/problem", contestController.HandleCreateProblem)
	adminGroup.PUT("/:contestid/:problemid", contestController.HandleUpdateProblem)
	adminGroup.DELETE("/:contestid/:problemid", contestController.HandleDeleteProblem)

	//Leaderboard/User Management
	adminGroup.PUT("/:contestid/leaderboard/:userid", contestController.HandleUpdateLeaderboardUser)
}
