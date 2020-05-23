package admin

import (
	"net/http"

	echoview "github.com/foolin/goview/supports/echoview-v4"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// IndexController implements web.Controller
type IndexController struct{}

// IndexControllerParams fx params struct
type IndexControllerParams struct {
	fx.In
}

// IndexControllerResult fx result struct
type IndexControllerResult struct {
	fx.Out

	IndexController *IndexController `name:"AdminIndexController"`
}

// ProvideIndexController to fx
func ProvideIndexController(p IndexControllerParams) IndexControllerResult {
	return IndexControllerResult{
		IndexController: &IndexController{},
	}
}

// InitRoutes wire up routes
func (c *IndexController) InitRoutes(g *echo.Group) {
	g.GET("/", get)
}

func get(ctx echo.Context) error {
	return echoview.Render(ctx, http.StatusOK, "index", echo.Map{})
}
