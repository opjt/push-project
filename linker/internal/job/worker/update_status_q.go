package worker

import (
	"fmt"
	"push/linker/internal/api/dto"
)

type UpdateStatusQueue struct {
	ch chan UpdateStatusJob
}
type UpdateStatusJob struct {
	DTO dto.UpdateMessageDTO
}

func NewUpdateStatusQueue() *UpdateStatusQueue {
	return &UpdateStatusQueue{
		ch: make(chan UpdateStatusJob, 10000),
	}
}

func (m *UpdateStatusQueue) Enqueue(dto dto.UpdateMessageDTO) error {
	select {
	case m.ch <- UpdateStatusJob{DTO: dto}:
		return nil
	default:
		return fmt.Errorf("update status queue full")
	}
}

func (m *UpdateStatusQueue) Channel() <-chan UpdateStatusJob {
	return m.ch
}

func (m *UpdateStatusQueue) Close() {
	close(m.ch)
}
