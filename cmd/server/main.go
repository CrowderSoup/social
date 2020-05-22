package main

import (
	"github.com/CrowderSoup/socialboat/internal/config"
	"github.com/CrowderSoup/socialboat/internal/controllers"
	"github.com/CrowderSoup/socialboat/internal/services"
	"github.com/CrowderSoup/socialboat/internal/web"

	"go.uber.org/fx"
)

func main() {
	bundle := fx.Options(
		config.Module,
		services.Module,
		web.Module,
		controllers.Module,
	)
	app := fx.New(
		bundle,
		fx.Invoke(services.InvokeDB),
		fx.Invoke(web.InvokeServer),
	)

	app.Run()

	<-app.Done()
}
