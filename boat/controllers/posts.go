package controllers

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/social/boat/models"
	"github.com/CrowderSoup/social/boat/services"
)

// PostsController controller for posts
type PostsController struct {
	DB *gorm.DB
}

// NewPostsController creates a new PostsController
func NewPostsController(db *gorm.DB) *PostsController {
	return &PostsController{
		DB: db,
	}
}

// InitRoutes initialize routes for this controller
func (c *PostsController) InitRoutes(g *echo.Group) {
	g.GET("", c.get)
	g.POST("", c.post)
}

func (c *PostsController) get(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	page, _ := strconv.Atoi(ctx.QueryParam("page"))
	limit, _ := strconv.Atoi(ctx.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}

	offset := 0
	if page > 1 {
		offset = (page - 1) * limit
	}

	var posts []models.Post
	c.DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&posts)

	return ctx.Render(http.StatusOK, "index", echo.Map{
		"title":    "SocialMast",
		"loggedIn": s.LoggedIn(),
		"posts":    posts,
	})
}

func (c *PostsController) post(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	title := ctx.FormValue("title")
	body := ctx.FormValue("body")

	if body == "" {
		panic("Body is required")
	}

	post := &models.Post{
		Title: title,
		Body:  body,
	}
	c.DB.Create(post)

	var posts []models.Post
	c.DB.Limit(10).Offset(0).Order("created_at desc").Find(&posts)

	return ctx.Render(http.StatusOK, "index", echo.Map{
		"title":    "SocialMast",
		"loggedIn": s.LoggedIn(),
		"posts":    posts,
	})
}
