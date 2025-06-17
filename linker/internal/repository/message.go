package repository

import (
	"context"
	"push/linker/internal/api/dto"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"
	"push/linker/types"
	"time"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *model.Message) (uint64, error)
	UpdateMessage(ctx context.Context, msg *model.Message) error
	UpdateMessages(ctx context.Context, dtos *dto.UpdateMessagesDTO) error
	GetById(ctx context.Context, id uint) (*model.Message, error)
	ReceiveMessage(ctx context.Context, msgId uint64) error
}
type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(mariaDb *database.MariaDB) MessageRepository {
	return &messageRepository{db: mariaDb.GetDB()}
}

func (r *messageRepository) GetById(ctx context.Context, id uint) (*model.Message, error) {
	var msg model.Message
	if err := r.db.WithContext(ctx).First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *messageRepository) CreateMessage(ctx context.Context, msg *model.Message) (uint64, error) {

	result := r.db.WithContext(ctx).Create(msg)

	return msg.ID, result.Error
}

func (r *messageRepository) UpdateMessage(ctx context.Context, msg *model.Message) error {

	updates := map[string]interface{}{
		"status": msg.Status,
	}

	if msg.SnsMsgId != "" {
		updates["sns_msg_id"] = msg.SnsMsgId
	}
	if msg.SentAt != nil {
		updates["sent_at"] = msg.SentAt
	}

	return r.db.Model(&model.Message{}).
		Where("id = ?", msg.ID).
		Updates(updates).Error

}
func (r *messageRepository) UpdateMessages(ctx context.Context, dto *dto.UpdateMessagesDTO) error {
	if len(dto.Ids) == 0 {
		return nil
	}

	updates := map[string]interface{}{
		"status": dto.Column.Status,
	}
	if dto.Column.SnsMsgId != "" {
		updates["sns_msg_id"] = dto.Column.SnsMsgId
	}

	return r.db.WithContext(ctx).
		Model(&model.Message{}).
		Where("id IN ?", dto.Ids).
		Updates(updates).Error
}

func (r *messageRepository) ReceiveMessage(ctx context.Context, msgId uint64) error {
	updates := map[string]interface{}{

		"status":  types.StatusSent,
		"sent_at": time.Now(),
	}
	return r.db.Model(&model.Message{}).
		Where("id = ? AND sent_at IS NULL", msgId).
		Updates(updates).Error
}
