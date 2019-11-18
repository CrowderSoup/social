package controllers

import (
	"net/http"
	"strings"

	"github.com/CrowderSoup/socialboat/services"
	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"
)

// ProfileController handles HTTP requests for the profile
type ProfileController struct {
	Service *services.ProfileService
}

// NewProfileController returns a new ProfileController
func NewProfileController(db *gorm.DB) *ProfileController {
	return &ProfileController{
		Service: services.NewProfileService(db),
	}
}

// InitRoutes initialize the routes for this controller
func (c *ProfileController) InitRoutes(g *echo.Group) {
	g.GET("", c.get)
	g.POST("", c.update)
}

func (c *ProfileController) get(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	profile, err := c.Service.GetByUserID(bc.Session.UserID())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get profile")
	}

	return bc.ReturnView(http.StatusOK, "profile", echo.Map{
		"profile": profile,
	})
}

func (c *ProfileController) update(ctx echo.Context) error {
	bc := ctx.(*BoatContext)

	profile, err := c.Service.GetByUserID(bc.Session.UserID())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get profile")
	}

	nickname := strings.TrimSpace(bc.FormValue("nickname"))
	firstName := strings.TrimSpace(bc.FormValue("first_name"))
	lastName := strings.TrimSpace(bc.FormValue("last_name"))
	publicEmail := strings.TrimSpace(bc.FormValue("public_email"))
	twitter := strings.TrimSpace(bc.FormValue("twitter"))
	github := strings.TrimSpace(bc.FormValue("github"))
	phone := strings.TrimSpace(bc.FormValue("phone"))
	note := strings.TrimSpace(bc.FormValue("note"))

	profile.NickName = nickname
	profile.FirstName = firstName
	profile.LastName = lastName
	profile.PublicEmail = publicEmail
	profile.Twitter = twitter
	profile.Github = github
	profile.Phone = phone
	profile.Note = note

	if err = c.Service.Update(profile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update profile")
	}

	return bc.ReturnView(http.StatusOK, "profile", echo.Map{
		"profile": profile,
	})
}
