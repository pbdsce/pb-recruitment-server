package main

import (
	"app/internal"
	"app/internal/boot"
	"app/internal/controllers"
	"app/internal/db"
	"app/internal/routes"
	"app/internal/services"
	"app/internal/stores"
	"context"
	"log"

	"go.uber.org/fx"
)

func main() {
	if err := boot.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	fx.New(
		fx.Provide(
			boot.NewFirebaseAuth,
			// Controllers
			controllers.NewContestController,
			controllers.NewUserController,
			controllers.NewSubmissionController,
			// Services
			services.NewContestService,
			services.NewUserService,
			services.NewSubmissionService,
			// Server
			internal.NewEchoServer,
			// Stores
			stores.NewStorage,
		),

		// Initialize database connection
		fx.Invoke(func() error {
			return db.InitDB()
		}),

		// Add routes to the Echo server
		fx.Invoke(routes.AddUserRoutes),
		fx.Invoke(routes.AddContestRoutes),
		fx.Invoke(routes.AddSubmissionRoutes),

		// Start the Echo server
		fx.Invoke(internal.StartEchoServer),

		// Graceful shutdown - close database connection
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return db.CloseDB()
				},
			})
		}),
	).Run()
}
