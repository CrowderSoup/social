package controllers

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/social/boat/models"
	"github.com/CrowderSoup/social/boat/services"
)

// AuthController auth controller
type AuthController struct {
	DB          *gorm.DB
	UserService *services.UserService
}

// NewAuthController creates a new AuthController
func NewAuthController(db *gorm.DB) *AuthController {
	userService := services.NewUserService(db)
	return &AuthController{
		DB:          db,
		UserService: userService,
	}
}

// InitRoutes initialize routes for AuthController
func (c *AuthController) InitRoutes(g *echo.Group) {
	g.GET("", c.get)
	g.GET("/logout", c.logout)
	g.POST("/login", c.login)
	g.POST("/register", c.register)
}

func (c *AuthController) get(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	if s.LoggedIn() {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	return ctx.Render(http.StatusOK, "auth", echo.Map{
		"title": "SocialMast - Auth",
	})
}

func (c *AuthController) login(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	if s.LoggedIn() {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	email := ctx.QueryParam("email")
	password := ctx.QueryParam("password")

	user, err := c.UserService.GetByEmail(email)
	if err != nil {
		return ctx.Render(http.StatusUnauthorized, "4xx", echo.Map{
			"title": "SocialMast - Uh Oh!",
			"error": fmt.Sprint(err),
		})
	}

	validPassword, err := c.UserService.CheckPassword(password, user)

	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "5xx", echo.Map{
			"title": "SocialMast - Oops!",
			"error": fmt.Sprint(err),
		})
	}

	if !validPassword {
		return ctx.Render(http.StatusUnauthorized, "4xx", echo.Map{
			"title": "SocialMast - Uh Oh!",
			"error": fmt.Sprint(err),
		})
	}

	err = s.SetValue("loggedIn", true, true)
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) register(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	if s.LoggedIn() {
		return ctx.Redirect(http.StatusSeeOther, "/")
	}

	email := ctx.QueryParam("email")
	password := ctx.QueryParam("password")

	user := &models.User{
		Email:    email,
		Password: password,
	}

	err = c.UserService.Create(user)
	if err != nil {
		return ctx.Render(http.StatusInternalServerError, "5xx", echo.Map{
			"title": "SocialMast - Opps!",
			"error": fmt.Sprint(err),
		})
	}

	err = s.SetValue("loggedIn", true, true)
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) logout(ctx echo.Context) error {
	s, err := services.GetSession("Boat", ctx)
	if err != nil {
		return err
	}

	err = s.ClearValue("loggedIn")
	if err != nil {
		return err
	}

	return ctx.Redirect(http.StatusSeeOther, "/")
}
