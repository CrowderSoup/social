package admin

import (
	"net/http"
	"strconv"

	"github.com/CrowderSoup/socialboat/internal/services"

	echo "github.com/labstack/echo/v4"
)

// PostsController implements web.Controller
type PostsController struct {
	PostService *services.PostService
}

// ProvidePostsController to fx
func ProvidePostsController(postService *services.PostService) *PostsController {
	return &PostsController{
		PostService: postService,
	}
}

// InitRoutes wire up routes
func (c *PostsController) InitRoutes(g *echo.Group) {
	g.GET("/posts", c.get)
}

func (c *PostsController) get(ctx echo.Context) error {
	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))

	posts, err := c.PostService.GetList(page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting posts")
	}

	return ctx.Render(http.StatusOK, "posts/index", echo.Map{
		"posts": posts,
	})
}
