package main

import (
	"app/internal"
	"app/internal/boot"
	"log"

	"go.uber.org/fx"
)

func main() {
	if err := boot.LoadEnv(); err != nil {
		log.Fatal(err)
	}

	fx.New(
		fx.Provide(boot.NewFirebaseAuth),

		// Start the Echo server
		fx.Invoke(internal.StartEchoServer),
	).Run()
}
