package core

import (
	"push/common"
	"push/sender/internal/pkg"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	pkg.Module,
	common.CommonModules,
)
