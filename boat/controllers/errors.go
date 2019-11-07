package controllers

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
)

// HTTPErrorHandler handle errors
func HTTPErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	message := "Something went wrong"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprint(he.Message)
	}

	// Log the error
	ctx.Logger().Error(err)

	page := "5xx"
	if code < http.StatusInternalServerError {
		page = "4xx"
	}

	fmt.Println(err)
	ctx.Render(http.StatusUnauthorized, page, echo.Map{
		"title": "SocialMast - Uh Oh!",
		"error": message,
		"code":  code,
	})
}
