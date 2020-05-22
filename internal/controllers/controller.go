package controllers

import (
	"go.uber.org/fx"

	"github.com/CrowderSoup/socialboat/internal/controllers/admin"
)

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		admin.ProvideIndexController,
	),
)
