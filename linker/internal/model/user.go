package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`

	Messages []Message `gorm:"foreignKey:UserID"` // has-many 관계
}
