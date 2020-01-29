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

	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

// PostsController controller for posts
type PostsController struct {
	PostService *services.PostService
}

// NewPostsController creates a new PostsController
func NewPostsController(db *gorm.DB) *PostsController {
	return &PostsController{
		PostService: services.NewPostService(db),
	}
}

// InitRoutes initialize routes for this controller
func (c *PostsController) InitRoutes(g *echo.Group) {
	g.GET("", c.listAll)
	g.GET("posts/:slug", c.singlePost)
	g.GET("posts/:slug/*", c.singlePost)
	g.GET("posts/:slug/edit", c.edit)
	g.POST("posts/:slug/update", c.update)
	g.POST("posts/:slug/delete", c.delete)
	g.POST("", c.create)
}

func (c *PostsController) listAll(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	page, _ := strconv.Atoi(bc.QueryParam("page"))
	limit, _ := strconv.Atoi(bc.QueryParam("limit"))

	posts, err := c.PostService.GetList(page, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error getting posts")
	}

	return bc.ReturnView(http.StatusOK, "index", echo.Map{
		"posts": posts,
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
		UserID: bc.Session.UserID(),
	}

	err = c.PostService.Create(post)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating post")
	}

	return bc.Redirect(http.StatusSeeOther, fmt.Sprintf("/posts/%s", URLSlug))
}

func (c *PostsController) edit(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	// Get post and pass to view
	slug := bc.Param("slug")

	post, err := c.PostService.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "post not found")
	}

	return bc.ReturnView(http.StatusOK, "post-edit", echo.Map{
		"title": fmt.Sprintf("Edit %s", post.Title),
		"post":  post,
	})
}

func (c *PostsController) update(ctx echo.Context) error {
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

	// Get the post
	slug := bc.Param("slug")
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

	return bc.Redirect(http.StatusSeeOther, fmt.Sprintf("/posts/%s", slug))
}

func (c *PostsController) delete(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	// Ensure the user is logged in
	err := bc.EnsureLoggedIn()
	if err != nil {
		return err
	}

	// Get the post
	slug := bc.Param("slug")
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

	return bc.Redirect(http.StatusSeeOther, "/")
}

func (c *PostsController) singlePost(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	slug := bc.Param("slug")
	if slug == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	post, err := c.PostService.GetBySlug(slug)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return bc.ReturnView(http.StatusOK, "post", echo.Map{
		"title": post.Title,
		"post":  post,
	})
}
