package controllers

import (
	"github.com/CrowderSoup/socialboat/internal/controllers/admin"

	echoview "github.com/foolin/goview/supports/echoview-v4"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// CustomContextMiddlewareResult middleware for injecting our custom context
type CustomContextMiddlewareResult struct {
	fx.Out

	CustomContextMiddleware echo.MiddlewareFunc `name:"CustomContextMiddleware"`
}

// ProvideCustomContextMiddleware returns our custom context middleware
func ProvideCustomContextMiddleware() CustomContextMiddlewareResult {
	return CustomContextMiddlewareResult{
		CustomContextMiddleware: func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				cc := &BoatContext{
					Context: c,
				}
				return next(cc)
			}
		},
	}
}

// BoatContext custom echo Context for boat
type BoatContext struct {
	echo.Context
}

// Render overrides echo's Render method and uses echoview.Render instead
// This ensures that the front and back end's get the renderer they need
func (bc *BoatContext) Render(status int, view string, data interface{}) error {
	return echoview.Render(bc.Context, status, view, data)
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		ProvideCustomContextMiddleware,
		admin.ProvideIndexController,
	),
)
