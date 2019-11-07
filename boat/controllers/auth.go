package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// AuthController auth controller
type AuthController struct {
	DB *gorm.DB
}

// NewAuthController creates a new AuthController
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		DB: db,
	}
}

// InitRoutes initialize routes for AuthController
func (c *AuthController) InitRoutes(g *echo.Group) {
	g.GET("", c.get)
	g.POST("/login", c.login)
	g.POST("/register", c.register)
}

func (c *AuthController) get(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "auth", echo.Map{
		"title": "SocialMast - Auth",
	})
}

func (c *AuthController) login(ctx echo.Context) error {
	return nil
}

func (c *AuthController) register(ctx echo.Context) error {
	return nil
}
