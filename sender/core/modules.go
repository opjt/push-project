package core

import (
	"push/common/lib"
	"push/sender/pkg/sqs"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	sqs.Module,
)
