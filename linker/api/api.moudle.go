package api

import (
	"push/linker/api/controller"
	"push/linker/api/router"

	"go.uber.org/fx"
)

var Module = fx.Options(
	router.Module,
	controller.Module,
)
