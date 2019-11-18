package controllers

import (
	"net/http"

	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

// AuthController auth controller
type AuthController struct {
	UserService    *services.UserService
	ProfileService *services.ProfileService
}

// NewAuthController creates a new AuthController
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		UserService:    services.NewUserService(db),
		ProfileService: services.NewProfileService(db),
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
	bc := ctx.(*BoatContext)
	bc.RedirectIfLoggedIn("/")

	return bc.ReturnView(http.StatusOK, "auth", echo.Map{
		"title": "SocialMast - Auth",
	})
}

func (c *AuthController) login(ctx echo.Context) error {
	bc := ctx.(*BoatContext)
	bc.RedirectIfLoggedIn("/")

	email := bc.FormValue("email")
	password := bc.FormValue("password")

	user, err := c.UserService.GetByEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	validPassword, err := c.UserService.CheckPassword(password, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	if !validPassword {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid email or password")
	}

	err = bc.Session.SetValue("loggedIn", true, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Problem setting session")
	}

	err = bc.Session.SetValue("userID", user.ID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Problem setting session")
	}

	return bc.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) register(ctx echo.Context) error {
	bc := ctx.(*BoatContext)
	bc.RedirectIfLoggedIn("/")

	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	user := &models.User{
		Email:    email,
		Password: password,
	}

	err := c.UserService.Create(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user")
	}

	// Create a Profile
	profile := &models.Profile{
		UserID: user.ID,
	}

	err = c.ProfileService.Create(profile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating user profile")
	}

	err = bc.Session.SetValue("loggedIn", true, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Problem setting session")
	}

	err = bc.Session.SetValue("userID", user.ID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Problem setting session")
	}

	return bc.Redirect(http.StatusSeeOther, "/")
}

func (c *AuthController) logout(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	err := bc.Session.ClearValue("loggedIn")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Problem clearing session")
	}

	return bc.Redirect(http.StatusSeeOther, "/")
}
