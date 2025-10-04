package main

import (
	"app/internal"
	"app/internal/boot"
	"app/internal/controllers"
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
			services.NewContestService,
			internal.NewEchoServer,
		),

		// Start the Echo server
		fx.Invoke(internal.StartEchoServer),
	).Run()
}
