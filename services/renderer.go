package services

import (
	"html/template"

	"github.com/foolin/goview"
	echoview "github.com/foolin/goview/supports/echoview-v4"
)

// RendererConfig configuration for our renderer
type RendererConfig struct {
	Root         string   `default:"views"`
	Extension    string   `default:".html"`
	Master       string   `default:"layouts/master"`
	Partials     []string `required:"true"`
	DisableCache bool     `default:"true"`
}

// NewRenderer build and return a new renderer
func NewRenderer(config RendererConfig) *echoview.ViewEngine {
	return echoview.New(goview.Config{
		Root:      config.Root,
		Extension: config.Extension,
		Master:    config.Master,
		Partials:  config.Partials,
		Funcs: template.FuncMap{
			"noescape": NoEscape,
		},
		DisableCache: config.DisableCache,
	})
}

// NoEscape returns an unescaped HTML string
func NoEscape(s string) template.HTML {
	return template.HTML(s)
}
