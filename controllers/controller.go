package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
	echo "github.com/labstack/echo/v4"

	"github.com/CrowderSoup/socialboat/config"
	"github.com/CrowderSoup/socialboat/models"
	"github.com/CrowderSoup/socialboat/services"
)

// CustomContextHandler struct with a func for handling the custom context intjection
type CustomContextHandler struct {
	Server         *config.Server
	ProfileService *services.ProfileService
	MenuService    *services.MenuService
}

// NewCustomContextHandler returns a new CustomContextHandler
func NewCustomContextHandler(db *gorm.DB, s *config.Server) *CustomContextHandler {
	return &CustomContextHandler{
		Server:         s,
		ProfileService: services.NewProfileService(db),
		MenuService:    services.NewMenuService(db),
	}
}

// Handler injects our custom context
func (h *CustomContextHandler) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := services.GetSession("Boat", c)
		if err != nil {
			return err
		}

		cc := &BoatContext{
			Context:        c,
			Server:         h.Server,
			Session:        s,
			ProfileService: h.ProfileService,
			MenuService:    h.MenuService,
		}
		return next(cc)
	}
}

// ManifestHandler handles the manifest
func ManifestHandler(ctx echo.Context) error {
	manifest, err := ioutil.ReadFile("./manifest.webmanifest")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error reading manifest")
	}

	return ctx.Blob(http.StatusOK, "application/manifest+json", manifest)
}

// BoatContext custom echo Context for boat
type BoatContext struct {
	echo.Context

	Server         *config.Server
	Session        *services.Session
	ProfileService *services.ProfileService
	MenuService    *services.MenuService
}

// LoggedIn checks if a contexts session is logged in
func (bc *BoatContext) LoggedIn() bool {
	return bc.Session.LoggedIn()
}

// EnsureLoggedIn ensures a user is logged in, throws error if not
func (bc *BoatContext) EnsureLoggedIn() error {
	if !bc.LoggedIn() {
		return echo.NewHTTPError(http.StatusUnauthorized, "You must be logged in")
	}

	return nil
}

// RedirectIfLoggedIn redirects to given path if logged in
func (bc *BoatContext) RedirectIfLoggedIn(path string) error {
	if bc.LoggedIn() {
		return bc.Redirect(http.StatusSeeOther, path)
	}

	return nil
}

// ReturnView renders a view, adding some data to the return
func (bc *BoatContext) ReturnView(code int, view string, data echo.Map) error {
	// Set "title" if not already set
	if _, ok := data["title"]; !ok {
		data["title"] = bc.Server.SiteName
	} else {
		data["title"] = fmt.Sprintf("%s: %s", bc.Server.SiteName, data["title"])
	}

	// Set "siteName" if not already set
	if _, ok := data["siteName"]; !ok {
		data["siteName"] = bc.Server.SiteName
	}

	// Set "tagline" if not already set
	if _, ok := data["tagline"]; !ok {
		data["tagline"] = bc.Server.TagLine
	}

	// Set "loggedIn" if not already set
	if _, ok := data["loggedIn"]; !ok {
		data["loggedIn"] = bc.LoggedIn()
	}

	// Set "profile" if not already set
	if _, ok := data["profile"]; !ok {
		profile, err := bc.ProfileService.GetFirst()

		// TODO: Clean up this hack, there has to be a better way...
		if err != nil && err.Error() != "record not found" {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't get profile")
		} else if err != nil && err.Error() == "record not found" && bc.Request().RequestURI != "/auth" {
			return bc.Redirect(http.StatusSeeOther, "/auth")
		}

		if profile == nil {
			profile = &models.Profile{
				NickName: "SocialMast",
			}
		}

		data["profile"] = profile
	}

	if _, ok := data["menus"]; !ok {
		menus, err := bc.MenuService.GetAllForView()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "couldn't load menus")
		}

		data["menus"] = menus
	}

	return bc.Render(code, view, data)
}
