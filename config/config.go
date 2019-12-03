package config

import (
	"github.com/CrowderSoup/socialboat/services"
	"github.com/koding/multiconfig"
)

// Server our server
type Server struct {
	SiteName       string `default:"SocialBoat"`
	TagLine        string `default:""`
	AssetsDir      string `default:"assets"`
	DBConfig       DBConfig
	Port           int `default:"8080"`
	RendererConfig services.RendererConfig
	SessionSecret  string `required:"true"`
	Migrate        bool   `default:"false"`
	MigrateUp      bool   `default:"true"`
}

// DBConfig Config for the database
type DBConfig struct {
	ConnectionString string `default:"boat.db"`
	Dialect          string `default:"sqlite3"`
}

// LoadConfig loads the config
func LoadConfig(s *Server, path string) {
	m := multiconfig.NewWithPath(path)
	m.MustLoad(s)
}
