package routes

import (
	"app/internal/controllers"
	"app/internal/middleware"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
)

func AddContestRoutes(
	e *echo.Echo,
	authClient *auth.Client,
	contestController *controllers.ContestController,
) {
	// List all contests
	e.GET("/contests/list",
		contestController.ListContests,
		middleware.OptionalFirebaseAuth(authClient),
	)

	// // Get details of a specific contest
	// // If the user is authenticated, return user-specific details
	// // If not, return public details
	// e.GET("/contests/:id",
	// 	contestController.GetContest,
	// 	middleware.OptionalFirebaseAuth(authClient),
	// )

	// // Get the leaderboard of a specific contest
	// // Paginate, page=<page> and 20 entries per page
	// e.GET("/contests/:id/leaderboard",
	// 	contestController.GetLeaderboard,
	// )
}
