package web

import (
	"context"
	"fmt"

	"github.com/CrowderSoup/socialboat/internal/config"

	echoview "github.com/foolin/goview/supports/echoview-v4"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wader/gormstore"
	"go.uber.org/fx"
)

// Server serves up our application
type Server struct {
	echo *echo.Echo
}

// NewWebServer returns a web server
func NewWebServer(
	e *echo.Echo,
	c *config.Config,
	r *echoview.ViewEngine,
) *Server {
	// Static dir
	e.Static(fmt.Sprintf("/%s", c.AssetsDir), c.AssetsDir)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Set Renderer
	e.Renderer = r

	return &Server{
		echo: e,
	}
}

// InvokeServer starts up our web server
func InvokeServer(
	lc fx.Lifecycle,
	c *config.Config,
	server *Server,
	sessionStore *gormstore.Store,
) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				// Server session store (wire up here because db is required)
				server.echo.Use(session.Middleware(sessionStore))

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
		ProvideAdminMiddleware,
		ProvideAdminGroup,
		NewWebServer,
	),
)
