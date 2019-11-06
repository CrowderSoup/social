package main

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/CrowderSoup/socialmast.xyz/controllers"
	"github.com/CrowderSoup/socialmast.xyz/models"
)

func main() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&models.Post{})

	r := gin.Default()

	r.HTMLRender = ginview.Default()

	store := cookie.NewStore([]byte("secret"))

	// Middleware
	r.Use(sessions.Sessions("socialmast", store))

	postsController := controllers.NewPostsController(db)

	r.GET("/", postsController.Get)
	r.POST("/", postsController.Post)

	r.Run() // listen and serve on :8080
}
