package repository

import (
	"context"
	"push/linker/internal/model"
	"push/linker/internal/pkg/database"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
}
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(mariaDb *database.MariaDB) UserRepository {
	return &userRepository{db: mariaDb.GetDB()}
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
