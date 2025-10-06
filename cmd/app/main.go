package main

import (
	"app/internal"
	"app/internal/boot"
	"app/internal/controllers"
	"app/internal/routes"
	"app/internal/services"
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
			controllers.NewContestController,
			controllers.NewUserController,
			services.NewContestService,
			services.NewUserService,
			internal.NewEchoServer,
		),

		// Add routes to the Echo server
		fx.Invoke(routes.AddUserRoutes),
		fx.Invoke(routes.AddContestRoutes),

		// Start the Echo server
		fx.Invoke(internal.StartEchoServer),
	).Run()
}
