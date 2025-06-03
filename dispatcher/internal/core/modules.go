package core

import (
	"push/common"
	"push/sender/internal/sender"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	sender.Module,
	common.CommonModules,
)
