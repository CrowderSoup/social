package controllers

import (
	"github.com/CrowderSoup/socialboat/internal/controllers/admin"

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

// ReturnView renders our view using the correct renderer
func (bc *BoatContext) ReturnView() error {
	panic("not implemnted")
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		ProvideCustomContextMiddleware,
		admin.ProvideIndexController,
	),
)
