package admin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/CrowderSoup/socialboat/internal/models"
	"github.com/CrowderSoup/socialboat/internal/services"

	"github.com/gosimple/slug"
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
	g.GET("/posts/new", c.new)
	g.GET("/posts/:slug", c.single)
	g.POST("/posts/new", c.saveNew)
	g.POST("/posts/:slug", c.update)
	g.POST("/posts/:slug/delete", c.delete)
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

func (c *PostsController) new(ctx echo.Context) error {
	post := models.Post{}

	return ctx.Render(http.StatusOK, "posts/edit", echo.Map{
		"title": "New Post",
		"post":  post,
	})
}

func (c *PostsController) single(ctx echo.Context) error {
	slug := ctx.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	post, err := c.PostService.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return ctx.Render(http.StatusOK, "posts/edit", echo.Map{
		"title": post.Title,
		"post":  post,
	})
}

func (c *PostsController) saveNew(ctx echo.Context) error {
	title := strings.TrimSpace(ctx.FormValue("title"))
	body := strings.TrimSpace(ctx.FormValue("body"))

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

	// If there's a post with the same slug, add timestamp to this posts slug
	existingPost, _ := c.PostService.GetBySlug(URLSlug)
	if existingPost != nil {
		timestamp := time.Now().Unix()
		URLSlug = slug.Make(fmt.Sprintf("%s %d", URLSlug, timestamp))
	}

	post := &models.Post{
		Title:  title,
		Body:   body,
		Slug:   URLSlug,
		UserID: 1, // TODO: Update to pull user id from session
	}

	if err := c.PostService.Create(post); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating post")
	}

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/posts/%s", post.Slug))
}

func (c *PostsController) update(ctx echo.Context) error {
	title := strings.TrimSpace(ctx.FormValue("title"))
	body := strings.TrimSpace(ctx.FormValue("body"))

	if body == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Body is required!")
	}

	// Get the post
	slug := ctx.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	post, err := c.PostService.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	// Update the post
	post.Title = title
	post.Body = body

	// Save the post
	err = c.PostService.Update(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update post")
	}

	return ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/posts/%s", post.Slug))
}

func (c *PostsController) delete(ctx echo.Context) error {
	// Get the post
	slug := ctx.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	post, err := c.PostService.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	err = c.PostService.Delete(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to delete")
	}

	return ctx.Redirect(http.StatusSeeOther, "/admin/posts")
}
