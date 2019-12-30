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
	RootURL        string `default:"http://localhost:8080"` // With Protocol, include port if your reverse proxy has your site on a non-standard port (80/443)
	Port           int    `default:"8080"`                  // Port the Go server will run on (enables you to run multiple instances on the same server)
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
