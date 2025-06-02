package repository

import (
	"context"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"

	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *model.Message) (uint, error)
	UpdateMessage(ctx context.Context, msg *model.Message) error
	GetById(ctx context.Context, id uint) (*model.Message, error)
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

func (r *messageRepository) CreateMessage(ctx context.Context, msg *model.Message) (uint, error) {

	result := r.db.WithContext(ctx).Create(msg)

	return msg.ID, result.Error
}

func (r *messageRepository) UpdateMessage(ctx context.Context, msg *model.Message) error {
	return r.db.Model(&model.Message{}).
		Where("id = ?", msg.ID).
		Updates(map[string]interface{}{
			"status":         msg.Status,
			"sns_message_id": msg.SnsMsgId,
		}).Error
}
