package model

import "time"

type Message struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    uint   `gorm:"not null;index"` // 외래키
	Content   string `gorm:"type:text"`
	Status    string `gorm:"type:varchar(20);not null;default:pending;index"`
	SentAt    *time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // 연관관계
}
