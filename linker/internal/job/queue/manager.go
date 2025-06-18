package queue

import (
	"push/linker/internal/api/dto"
)

type JobQueueManager struct {
	updateStatusQ *UpdateStatusQueue
}

func NewJobQueueManager(updateStatusQ *UpdateStatusQueue) *JobQueueManager {
	return &JobQueueManager{
		updateStatusQ: updateStatusQ,
	}
}

func (m *JobQueueManager) EnqUpdateStatus(dto dto.UpdateMessageDTO) error {
	// 향후 로깅, 트레이싱, 중복검사, 통계 등 확장 가능
	return m.updateStatusQ.enqueue(dto)
}
