package model

import "time"

type Message struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement"`
	UserID    uint64     `gorm:"not null;index"` // 외래키
	Title     string     `gorm:"type:text"`
	Content   string     `gorm:"type:text"`
	SnsMsgId  string     `gorm:"type:varchar(100);index"`
	Status    string     `gorm:"type:varchar(20);not null;default:pending;index"`
	SentAt    *time.Time //nil 처리를 위해 포인터
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"` // 연관관계
}
