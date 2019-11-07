package controllers

import (
	session "github.com/ipfans/echo-session"
	"github.com/labstack/echo"
)

// Controller is an interface for controllers
type Controller interface {
}

// GetSessionValue gets a session value
func GetSessionValue(ctx echo.Context, key string) interface{} {
	s := session.Default(ctx)
	return s.Get(key)
}

// SetSessionValue sets a session value
func SetSessionValue(ctx echo.Context, key, value string) {
	s := session.Default(ctx)
	s.Set(key, value)
	s.Save()
}

// ClearSessionValue clears the current session
func ClearSessionValue(ctx echo.Context, key string) {
	s := session.Default(ctx)
	s.Delete(key)
	s.Save()
}
