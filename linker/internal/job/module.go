package job

import (
	"push/linker/internal/job/manager"
	"push/linker/internal/job/queue"
	"push/linker/internal/job/worker"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(queue.NewUpdateStatusQueue),
	fx.Provide(manager.NewJobQueueManager),
	fx.Provide(worker.NewJobUpdateStatus),
	fx.Invoke(func(j *worker.JobUpdateStatus) {}),
)
