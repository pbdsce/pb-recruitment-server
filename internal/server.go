package internal

import (
	"app/internal/controllers"
	"app/internal/middleware"
	"context"

	"firebase.google.com/go/v4/auth"
	"github.com/labstack/echo/v4"
	mdw "github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

func NewEchoServer(
	authClient *auth.Client,
	contestController *controllers.ContestController,
) *echo.Echo {
	e := echo.New()
	e.Use(mdw.Recover())
	e.Use(mdw.Logger())
	e.Use(mdw.CORS())

	e.GET("/contests/list",
		contestController.ListContests,
		middleware.OptionalFirebaseAuth(authClient),
	)

	e.POST("/contests/:id/register",
		contestController.RegisterParticipant,
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
