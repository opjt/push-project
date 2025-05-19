package bootstrap

import (
	"push/common/lib"
	"push/linker/internal/api"
	"push/linker/internal/api/service"
	"push/linker/internal/core"
	"push/linker/internal/pkg"
	"push/linker/internal/repository"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	lib.Module,
	pkg.Module,
	service.Module,
	repository.Module,
	fx.Provide(core.NewEngine),
)
