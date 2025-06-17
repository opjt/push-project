package worker

import (
	"context"
	"fmt"
	"push/common/lib/logger"
	"push/linker/internal/api/dto"
	"push/linker/internal/service"
	"sync"
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
	wg      sync.WaitGroup
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
			j.wg.Add(1) // 워커 1개
			go j.startProcessor()
			return nil
		},
		OnStop: func(context.Context) error {
			close(j.ch)
			j.wg.Wait()
			cancel()
			return nil
		},
	})
	return j
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
	defer j.wg.Done()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	var batch []UpdateStatusJob

	for {
		select {

		case job, ok := <-j.ch:
			if !ok {
				j.logger.Info("JobUpdateStatus processor stopped")
				j.flush(batch)
				return
			}
			batch = append(batch, job)

		case <-ticker.C:
			if len(batch) == 0 {
				continue
			}

			j.flush(batch)
			batch = nil
		}
	}
}

func (j *JobUpdateStatus) flush(batch []UpdateStatusJob) {
	if len(batch) == 0 {
		return
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
	ctx, cancel := context.WithTimeout(j.ctx, 5*time.Second)
	defer cancel()
	// 2. 그룹별 처리
	for key, ids := range groupMap {
		if err := j.service.UpdateMessagesStatus(ctx, ids, key); err != nil {
			j.logger.Errorf("Failed to batch update message status", err)
		}
	}
}
