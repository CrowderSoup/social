package config

import (
	"github.com/CrowderSoup/socialboat/services"
	"github.com/koding/multiconfig"
	"go.uber.org/fx"
)

// Config our server
type Config struct {
	SiteName       string `default:"SocialBoat"`
	TagLine        string `default:""`
	AssetsDir      string `default:"assets"`
	DBConfig       DBConfig
	RootURL        string `default:"http://localhost:8080"`
	Port           int    `default:"8080"`
	RendererConfig services.RendererConfig
	SessionSecret  string `required:"true"`
}

// DBConfig Config for the database
type DBConfig struct {
	ConnectionString string `default:"boat.db"`
	Dialect          string `default:"sqlite3"`
}

// ProvideConfig provides configuration for our app
func ProvideConfig() *Config {
	var config Config
	m := multiconfig.NewWithPath("config.toml")
	m.MustLoad(&config)

	return &config
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		ProvideConfig,
	),
)
