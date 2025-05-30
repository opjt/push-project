package core

import (
	"push/common"
	"push/linker/internal/api"
	"push/linker/internal/api/service"
	"push/linker/internal/pkg"
	"push/linker/internal/pkg/gin"
	"push/linker/internal/repository"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	api.Module,
	common.CommonModules,
	pkg.Module,
	service.Module,
	repository.Module,
	fx.Provide(gin.NewEngine),
)
