package common

import (
	"push/common/lib"
	"push/common/pkg"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(

	lib.Module,
	pkg.Module,
)
