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

// Group a group of routes
type Group interface {
	InitControllers(*echo.Echo)
}

// Controller a controller composed of handlers
type Controller interface {
	InitRoutes(*echo.Group)
}

// NewServerParams params for new server
type NewServerParams struct {
	fx.In

	Instance   *echo.Echo
	Config     *config.Config
	Renderer   *echoview.ViewEngine
	AdminGroup Group `name:"AdminGroup"`
}

// NewWebServer returns a web server
func NewWebServer(p NewServerParams) *Server {
	// Static dir
	p.Instance.Static(fmt.Sprintf("/%s", p.Config.AssetsDir), p.Config.AssetsDir)

	// Middleware
	p.Instance.Use(middleware.Logger())
	p.Instance.Use(middleware.Recover())

	//Set Renderer
	p.Instance.Renderer = p.Renderer

	// Init Routes
	p.AdminGroup.InitControllers(p.Instance)

	return &Server{
		echo: p.Instance,
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
