package core

import (
	"push/common"
	"push/linker/internal/api"
	"push/linker/internal/pkg"
	"push/linker/internal/pkg/gin"
	"push/linker/internal/repository"
	"push/linker/internal/service"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	common.CommonModules,
	pkg.Module,
	repository.Module,
	service.Module,
	fx.Provide(gin.NewEngine),
)
