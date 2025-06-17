package worker

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewJobUpdateStatus),
	fx.Invoke(func(j *JobUpdateStatus) {}),
)
