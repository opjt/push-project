package core

import (
	"push/common"
	"push/linker/internal/api"
	"push/linker/internal/job"
	"push/linker/internal/pkg"
	"push/linker/internal/pkg/gin"
	"push/linker/internal/repository"
	"push/linker/internal/service"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	common.CommonModules,
	api.Module,
	repository.Module,
	service.Module,
	job.Module,
	pkg.Module,
	fx.Provide(gin.NewEngine),
)
