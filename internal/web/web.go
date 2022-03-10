package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

// NewServerParams params for new server
type NewServerParams struct {
	fx.In

	Instance                *echo.Echo
	Config                  *config.Config
	Renderer                *echoview.ViewEngine
	AdminGroup              *AdminGroup         `name:"AdminGroup"`
	CustomContextMiddleware echo.MiddlewareFunc `name:"CustomContextMiddleware"`
}

// NewWebServer returns a web server
func NewWebServer(p NewServerParams) *Server {
	// Remove Trailing Slashes from requests
	p.Instance.Pre(middleware.RemoveTrailingSlash())

	// Static dir
	p.Instance.Static(fmt.Sprintf("/%s", p.Config.AssetsDir), p.Config.AssetsDir)

	// Middleware
	p.Instance.Use(middleware.Logger())
	p.Instance.Use(middleware.Recover())

	//Set Renderer
	p.Instance.Renderer = p.Renderer

	// Custom Context
	p.Instance.Use(p.CustomContextMiddleware)

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

				data, err := json.MarshalIndent(server.echo.Routes(), "", "  ")
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(string(data))

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
