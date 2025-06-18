package manager

import (
	"push/linker/internal/api/dto"
	"push/linker/internal/job/worker"
)

type JobQueueManager struct {
	updateStatusQ *worker.UpdateStatusQueue
}

func NewJobQueueManager(updateStatusQ *worker.UpdateStatusQueue) *JobQueueManager {
	return &JobQueueManager{
		updateStatusQ: updateStatusQ,
	}
}

func (m *JobQueueManager) Enqueue(dto dto.UpdateMessageDTO) error {
	// 향후 로깅, 트레이싱, 중복검사, 통계 등 확장 가능
	return m.updateStatusQ.Enqueue(dto)
}
