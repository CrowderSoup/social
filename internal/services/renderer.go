package services

import (
	"bytes"
	"html/template"

	"github.com/CrowderSoup/socialboat/internal/config"
	echo "github.com/labstack/echo/v4"
	"go.uber.org/fx"

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
func NewRenderer(config *config.Config) *echoview.ViewEngine {
	rendererCfg := config.RendererConfig

	return echoview.New(goview.Config{
		Root:      rendererCfg.Root,
		Extension: rendererCfg.Extension,
		Master:    rendererCfg.Master,
		Partials:  rendererCfg.Partials,
		Funcs: template.FuncMap{
			"markdown": Markdown,
		},
		DisableCache: rendererCfg.DisableCache,
	})
}

// BackendRendererResult our result struct for holding the backend render middleware
type BackendRendererResult struct {
	fx.Out

	RendererMiddleware echo.MiddlewareFunc `name:"BackendRendererMiddleware"`
}

// NewBackendRenderer returns a middleware for rendering backend
func NewBackendRenderer() BackendRendererResult {
	m := echoview.NewMiddleware(goview.Config{
		Root:         "views/admin",
		Extension:    ".html",
		Master:       "layouts/master",
		DisableCache: true,
	})

	return BackendRendererResult{
		RendererMiddleware: m,
	}
}

// Markdown returns rendered markdown
func Markdown(s string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(s), &buf); err != nil {
		panic(err)
	}

	return template.HTML(buf.String())
}
