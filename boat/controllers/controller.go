package controllers

import (
	"net/http"

	"github.com/CrowderSoup/social/boat/services"
	echo "github.com/labstack/echo/v4"
)

// BoatContext custom echo Context for boat
type BoatContext struct {
	echo.Context

	Session *services.Session
}

// LoggedIn checks if a contexts session is logged in
func (bc *BoatContext) LoggedIn() bool {
	return bc.Session.LoggedIn()
}

// RedirectIfLoggedIn redirects to given path if logged in
func (bc *BoatContext) RedirectIfLoggedIn(path string) error {
	if bc.LoggedIn() {
		return bc.Redirect(http.StatusSeeOther, path)
	}

	return nil
}
