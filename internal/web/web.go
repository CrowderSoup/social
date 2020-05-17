package web

import (
	"context"
	"fmt"

	"github.com/CrowderSoup/socialboat/internal/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
)

// Server serves up our application
type Server struct {
	echo *echo.Echo
}

// NewWebServer returns a web server
func NewWebServer(e *echo.Echo, c *config.Config) *Server {
	// Static dir
	e.Static(fmt.Sprintf("/%s", c.AssetsDir), c.AssetsDir)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Server{
		echo: e,
	}
}

// InvokeServer starts up our web server
func InvokeServer(lc fx.Lifecycle, c *config.Config, server *Server) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go server.echo.Start(fmt.Sprintf(":%d", c.Port))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.echo.Close()
			},
		},
	)
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		echo.New,
		NewWebServer,
	),
)
