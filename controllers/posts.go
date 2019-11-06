package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"github.com/CrowderSoup/socialmast.xyz/models"
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

// Get gets posts
func (c *PostsController) Get(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset := 0

	if page > 1 {
		offset = page * limit
	}

	var posts []models.Post
	c.DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&posts)

	ctx.HTML(http.StatusOK, "index", gin.H{
		"title": "SocialMast",
		"posts": posts,
	})
}

// Post save a post
func (c *PostsController) Post(ctx *gin.Context) {
	title := ctx.PostForm("title")
	body := ctx.PostForm("body")

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

	ctx.HTML(http.StatusOK, "index", gin.H{
		"title": "SocialMast",
		"posts": posts,
	})
}
