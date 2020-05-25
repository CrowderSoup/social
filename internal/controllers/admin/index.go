package admin

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// IndexController implements web.Controller
type IndexController struct{}

// ProvideIndexController to fx
func ProvideIndexController() *IndexController {
	return &IndexController{}
}

// InitRoutes wire up routes
func (c *IndexController) InitRoutes(g *echo.Group) {
	g.GET("/", c.get)
}

func (c *IndexController) get(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index", echo.Map{})
}
