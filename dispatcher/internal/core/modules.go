package core

import (
	"push/common"
	"push/dispatcher/internal/sender"
	"push/dispatcher/internal/sessionmanager"

	"go.uber.org/fx"
)

var Modules = fx.Options(
	sender.Module,
	sessionmanager.Module,
	common.CommonModules,
)
