package services

import (
	"bytes"
	"html/template"

	"github.com/foolin/goview"
	echoview "github.com/foolin/goview/supports/echoview-v4"
	"github.com/yuin/goldmark"
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
			"markdown": Markdown,
		},
		DisableCache: config.DisableCache,
	})
}

// Markdown returns rendered markdown
func Markdown(s string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(s), &buf); err != nil {
		panic(err)
	}

	return template.HTML(buf.String())
}
