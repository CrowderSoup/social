package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/social/boat/models"
	"github.com/CrowderSoup/social/boat/services"
)

// PostsController controller for posts
type PostsController struct {
	Service *services.PostService
}

// NewPostsController creates a new PostsController
func NewPostsController(db *gorm.DB) *PostsController {
	return &PostsController{
		Service: services.NewPostService(db),
	}
}

// InitRoutes initialize routes for this controller
func (c *PostsController) InitRoutes(g *echo.Group) {
	g.GET("", c.get)
	g.POST("", c.post)
}

func (c *PostsController) get(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	page, _ := strconv.Atoi(bc.QueryParam("page"))
	limit, _ := strconv.Atoi(bc.QueryParam("limit"))

	posts, err := c.Service.GetList(page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting posts")
	}

	return bc.Render(http.StatusOK, "index", echo.Map{
		"title":    "SocialMast",
		"loggedIn": bc.LoggedIn(),
		"posts":    posts,
	})
}

func (c *PostsController) post(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	title := strings.TrimSpace(bc.FormValue("title"))
	body := strings.TrimSpace(bc.FormValue("body"))

	if body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Body is required!")
	}

	var URLSlug string
	if title != "" {
		URLSlug = slug.Make(title)
	} else {
		URLSlug = slug.Make(body[:50])
	}

	post := &models.Post{
		Title: title,
		Body:  body,
		Slug:  URLSlug,
	}

	err := c.Service.Create(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting posts")
	}

	posts, err := c.Service.GetList(1, 10)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting posts")
	}

	return bc.Render(http.StatusOK, "index", echo.Map{
		"title":    "SocialMast",
		"loggedIn": bc.LoggedIn(),
		"posts":    posts,
	})
}
