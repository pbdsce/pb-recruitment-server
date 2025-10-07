package routes

import (
	"app/internal/controllers"

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
		//SAMPLE FETCH ROUTE
		contestController.ListContests,
		//middleware.OptionalFirebaseAuth(authClient),
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

	// // Register/Unregister the authenticated user for a specific contest
	// // Use a query parameter action=register or action=unregister
	// e.POST("/contests/:id/registration",
	// 	contestController.ModifyRegistration,
	// 	middleware.RequireFirebaseAuth(authClient),
	// )

	// // Get the problems of a specific contest for the authenticated user
	// // Do not return the problem statements themselves
	// e.GET("/contests/:id/problems",
	// 	contestController.GetContestProblemsList,
	// 	middleware.RequireFirebaseAuth(authClient),
	// )

	// // Get the problem statement of a specific problem in a contest for the authenticated user
	// e.GET("/contests/:id/problems/:problem_id",
	// 	contestController.GetContestProblemStatement,
	// 	middleware.RequireFirebaseAuth(authClient),
	// )
}
