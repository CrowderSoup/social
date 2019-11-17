package services

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo-contrib/session"
	echo "github.com/labstack/echo/v4"
	"github.com/wader/gormstore"
)

// Session a session that holds an internal session storage mechanism
type Session struct {
	Name     string
	Internal *sessions.Session
	Context  echo.Context
}

// InitSessionStore initializes the session storeage mechanism
func InitSessionStore(secret string, db *gorm.DB, shouldCleanup bool) *gormstore.Store {
	store := gormstore.New(db, []byte(secret))

	if shouldCleanup {
		quit := make(chan struct{})
		go store.PeriodicCleanup(1*time.Hour, quit)
	}

	return store
}

// GetSession gets a session
func GetSession(name string, ctx echo.Context) (*Session, error) {
	s := &Session{
		Name:    name,
		Context: ctx,
	}

	sess, err := session.Get(name, ctx)
	if err != nil {
		return nil, err
	}

	s.Internal = sess

	return s, nil
}

// GetValue gets a value from the session
func (s *Session) GetValue(key string) interface{} {
	return s.Internal.Values[key]
}

// SetValue sets a value in the session
func (s *Session) SetValue(key, value interface{}, shouldSave bool) error {
	s.Internal.Values[key] = value

	if shouldSave {
		return s.Save()
	}

	return nil
}

// ClearAll clears all values from the session
func (s *Session) ClearAll() error {
	for k := range s.Internal.Values {
		delete(s.Internal.Values, k)
	}

	return s.Save()
}

// ClearValue clears a single value from the session
func (s *Session) ClearValue(key interface{}) error {
	delete(s.Internal.Values, key)

	return s.Save()
}

// Save saves a session
func (s *Session) Save() error {
	err := s.Internal.Save(s.Context.Request(), s.Context.Response())
	if err != nil {
		return err
	}

	return nil
}

// LoggedIn checks the session to see if the user is logged in
func (s *Session) LoggedIn() bool {
	v := s.GetValue("loggedIn")
	loggedIn := false
	if l, ok := v.(bool); ok {
		loggedIn = l
	}

	return loggedIn
}

// UserID attempts to get the userID from session
func (s *Session) UserID() int {
	v := s.GetValue("userID")
	userID := 0
	if uid, ok := v.(uint); ok {
		userID = int(uid)
	}

	return userID
}
