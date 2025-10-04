package internal

import (
	"app/internal/middleware"
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func createEchoRoutes(e *echo.Echo, authClient *auth.Client) {
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	}, middleware.FirebaseAuth(authClient, true))
}

func StartEchoServer(lc fx.Lifecycle, authClient *auth.Client) *echo.Echo {
	e := echo.New()
	createEchoRoutes(e, authClient)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go e.Start(":8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
