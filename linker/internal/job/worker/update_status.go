package worker

import (
	"context"
	"push/common/lib/logger"
	"push/linker/internal/api/dto"
	"push/linker/internal/job/queue"
	"push/linker/internal/service"
	"push/linker/types"
	"sync"
	"time"

	"go.uber.org/fx"
)

type JobUpdateStatus struct {
	ch      <-chan queue.UpdateStatusJob
	ctx     context.Context
	service service.MessageService
	wg      sync.WaitGroup
	logger  *logger.Logger
}

func NewJobUpdateStatus(lc fx.Lifecycle, service service.MessageService, logger *logger.Logger, queue *queue.UpdateStatusQueue) *JobUpdateStatus {
	ctx, cancel := context.WithCancel(context.Background())
	j := &JobUpdateStatus{
		ch:      queue.Channel(),
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
			queue.Close()
			j.wg.Wait()
			cancel()
			return nil
		},
	})
	return j
}

func (j *JobUpdateStatus) startProcessor() {
	defer j.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var batch []queue.UpdateStatusJob

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

func (j *JobUpdateStatus) flush(batch []queue.UpdateStatusJob) {
	if len(batch) == 0 {
		return
	}

	// ID별 최신 상태만 유지
	type grouped struct {
		Status   string
		SnsMsgId string
	}
	idMap := make(map[uint64]grouped)

	for _, job := range batch {
		id := job.DTO.Id
		newStatus := job.DTO.Status
		newSnsMsgId := job.DTO.SnsMsgId

		prev, exists := idMap[id]
		if !exists {
			idMap[id] = grouped{Status: newStatus, SnsMsgId: newSnsMsgId}
			continue
		}

		// sent가 있으면 sent 우선
		if newStatus == types.StatusSent {
			// 기존에 sending이면 snsmsgid 유지
			if prev.SnsMsgId != "" && newSnsMsgId == "" {
				newSnsMsgId = prev.SnsMsgId
			}
			idMap[id] = grouped{Status: types.StatusSent, SnsMsgId: newSnsMsgId}
		} else if newStatus == types.StatusDeferred {
			idMap[id] = grouped{Status: types.StatusDeferred, SnsMsgId: prev.SnsMsgId}
		}
	}

	// 그룹핑: (status, snsmsgid) 조합으로
	groupMap := make(map[dto.UpdateMessageField][]uint64)
	for id, info := range idMap {
		key := dto.UpdateMessageField{
			Status:   info.Status,
			SnsMsgId: info.SnsMsgId,
		}
		groupMap[key] = append(groupMap[key], id)
	}

	// 전송
	ctx, cancel := context.WithTimeout(j.ctx, 5*time.Second)
	defer cancel()

	for key, ids := range groupMap {
		if err := j.service.UpdateMessagesStatus(ctx, ids, key); err != nil {
			j.logger.Errorf("Failed to batch update message status: %v", err)
		}
	}
}
