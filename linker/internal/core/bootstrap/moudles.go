package bootstrap

import (
	"push/common/lib"
	"push/linker/internal/api"
	"push/linker/internal/core"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	lib.Module,
	fx.Provide(core.NewEngine),
)
