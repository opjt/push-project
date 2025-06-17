package worker

import (
	"context"
	"fmt"
	"push/common/lib/logger"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"
	"time"

	"go.uber.org/fx"
)

type UpdateStatusJob struct {
	DTO dto.UpdateMessageDTO
}

type JobUpdateStatus struct {
	ch      chan UpdateStatusJob
	ctx     context.Context
	service service.MessageService
	logger  *logger.Logger
}

func NewJobUpdateStatus(lc fx.Lifecycle, service service.MessageService, logger *logger.Logger) *JobUpdateStatus {
	ctx, cancel := context.WithCancel(context.Background())
	j := &JobUpdateStatus{
		ch:      make(chan UpdateStatusJob, 1000),
		ctx:     ctx,
		service: service,
		logger:  logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go j.startProcessor()
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			close(j.ch)
			return nil
		},
	})
	return j
}

func RegisterJobUpdateStatus(lc fx.Lifecycle, service service.MessageService, logger *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			NewJobUpdateStatus(lc, service, logger)
			return nil
		},
	})
}

func (j *JobUpdateStatus) Enqueue(dto dto.UpdateMessageDTO) error {
	job := UpdateStatusJob{DTO: dto}
	select {
	case j.ch <- job:

		return nil
	default:
		// 채널이 가득 찼을 경우 개선 필요
		return fmt.Errorf("UpdateStatusQueue is full, dropping job")
	}
}

func (j *JobUpdateStatus) startProcessor() {

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	var batch []UpdateStatusJob

	for {
		select {
		case <-j.ctx.Done():
			return

		case job := <-j.ch:
			batch = append(batch, job)

		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}

			// 1. 그룹핑
			groupMap := make(map[dto.UpdateMessageField][]uint64)

			for _, job := range batch {
				key := dto.UpdateMessageField{
					Status:   job.DTO.Status,
					SnsMsgId: job.DTO.SnsMsgId,
				}
				groupMap[key] = append(groupMap[key], job.DTO.Id)
			}

			// 2. 그룹별로 처리
			for key, ids := range groupMap {

				if err := j.service.UpdateMessagesStatus(j.ctx, ids, key); err != nil {
					j.logger.Errorf("Failed to batch update message status", err)
				}
			}

			batch = nil
		}
	}
}
