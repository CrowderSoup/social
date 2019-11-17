package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/CrowderSoup/social/boat/services"
	echo "github.com/labstack/echo/v4"
)

// CustomContextHandler injects our custom context
func CustomContextHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		s, err := services.GetSession("Boat", c)
		if err != nil {
			return err
		}

		cc := &BoatContext{
			Context: c,
			Session: s,
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

	Session *services.Session
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
