package pkg

import (
	"push/common/pkg/aws"

	"go.uber.org/fx"
)

// Lib 모듈관리.
var Module = fx.Options(

	fx.Provide(aws.NewAwsConfig),
)
