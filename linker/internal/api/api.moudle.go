package api

import (
	"push/linker/internal/api/controller"
	"push/linker/internal/api/router"

	"go.uber.org/fx"
)

var Module = fx.Options(
	router.Module,
	controller.Module,
)
