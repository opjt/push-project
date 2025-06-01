package pkg

import (
	"push/common/pkg/awsinfra"

	"go.uber.org/fx"
)

// Lib 모듈관리.
var Module = fx.Options(

	fx.Provide(awsinfra.NewAwsConfig),
)
