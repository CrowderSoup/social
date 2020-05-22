package web

import (
	echo "github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// AdminGroupParams params for ProvideAdminGroup
type AdminGroupParams struct {
	fx.In

	Server     *Server
	Middleware []echo.MiddlewareFunc `name:"AdminMiddleware"`
}

// AdminGroupResult fx result struct for ProvideAdminGroup
type AdminGroupResult struct {
	fx.Out

	AdminGroup *echo.Group `name:"AdminGroup"`
}

// ProvideAdminMiddlewareParams fx params struct
type ProvideAdminMiddlewareParams struct {
	fx.In

	BackendRendererMiddleware echo.MiddlewareFunc `name:"BackendRendererMiddleware"`
}

// ProvideAdminMiddlewareResult fx result struct
type ProvideAdminMiddlewareResult struct {
	fx.Out

	AdminMiddleware []echo.MiddlewareFunc `name:"AdminMiddleware"`
}

// ProvideAdminMiddleware provides middleware for admin
func ProvideAdminMiddleware(p ProvideAdminMiddlewareParams) ProvideAdminMiddlewareResult {
	return ProvideAdminMiddlewareResult{
		AdminMiddleware: []echo.MiddlewareFunc{
			p.BackendRendererMiddleware,
		},
	}
}

// ProvideAdminGroup provides the admin group
func ProvideAdminGroup(p AdminGroupParams) AdminGroupResult {
	server := p.Server

	// Create echo group
	g := server.echo.Group("/admin", p.Middleware...)

	// TODO: Wire up group routes

	return AdminGroupResult{
		AdminGroup: g,
	}
}
