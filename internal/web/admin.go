package web

import (
	echo "github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

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

// AdminGroupParams params for ProvideAdminGroup
type AdminGroupParams struct {
	fx.In

	Instance   *echo.Echo
	Middleware []echo.MiddlewareFunc `name:"AdminMiddleware"`

	// Controllers
	IndexController Controller `name:"AdminIndexController"`
}

// AdminGroupResult fx result struct for ProvideAdminGroup
type AdminGroupResult struct {
	fx.Out

	AdminGroup Group `name:"AdminGroup"`
}

// AdminGroup group for admin routes
type AdminGroup struct {
	middleware []echo.MiddlewareFunc

	// Controllers
	indexController Controller
}

// ProvideAdminGroup provides the admin group
func ProvideAdminGroup(p AdminGroupParams) AdminGroupResult {
	return AdminGroupResult{
		AdminGroup: &AdminGroup{
			middleware:      p.Middleware,
			indexController: p.IndexController,
		},
	}
}

// InitControllers wire's up all our controllers
func (g *AdminGroup) InitControllers(instance *echo.Echo) {
	// Create echo group
	group := instance.Group("/admin", g.middleware...)

	// Wire up Controllers
	g.indexController.InitRoutes(group)
}
