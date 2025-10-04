package internal

import (
	"app/internal/controllers"
	"app/internal/middleware"
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewEchoServer(
	authClient *auth.Client,
	contestController *controllers.ContestController,
) *echo.Echo {
	e := echo.New()

	e.GET("/contests/list",
		contestController.ListContests,
		middleware.RequireFirebaseAuth(authClient),
	)

	return e
}

func StartEchoServer(lc fx.Lifecycle, e *echo.Echo) *echo.Echo {
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
