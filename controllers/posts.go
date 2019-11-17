package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	g.GET("", c.listAll)
	g.GET("posts/:slug", c.singlePost)
	g.GET("posts/:slug/*", c.singlePost)
	g.POST("", c.create)
}

func (c *PostsController) listAll(ctx echo.Context) error {
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

func (c *PostsController) create(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	title := strings.TrimSpace(bc.FormValue("title"))
	body := strings.TrimSpace(bc.FormValue("body"))

	if body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Body is required!")
	}

	var slugSource string
	if title != "" {
		slugSource = title
	} else {
		slugSource = body
	}

	// Ensure slug is the right length
	if len(slugSource) > 50 {
		slugSource = slugSource[:50]
	}

	URLSlug := slug.Make(slugSource)

	// We don't care about errors here probably
	existingPost, _ := c.Service.GetBySlug(URLSlug)
	if existingPost != nil {
		timestamp := time.Now().Unix()
		URLSlug = slug.Make(fmt.Sprintf("%s %d", URLSlug, timestamp))
	}

	post := &models.Post{
		Title:  title,
		Body:   body,
		Slug:   URLSlug,
		UserID: bc.Session.UserID(),
	}

	err = c.Service.Create(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating post")
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

func (c *PostsController) singlePost(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	slug := bc.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	post, err := c.Service.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get post")
	}

	return bc.Render(http.StatusOK, "post", echo.Map{
		"title":    "SocialMast - " + post.Title,
		"loggedIn": bc.LoggedIn(),
		"post":     post,
	})
}
