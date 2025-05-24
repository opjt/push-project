package model

import "time"

const (
	STATUS_PENDING = "pending" // 메시지 발행 후, 아직 sender가 처리 전인 상태
	STATUS_SENDING = "sending" // sender가 메시지를 받아 실제 전송 작업 중인 상태
	STATUS_SENT    = "sent"    // 메시지 전송 완료 상태, client에서 확인된 경우.
	STATUS_FAILED  = "failed"  // 메시지 전송 실패 상태
)

type Message struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	UserID    uint       `gorm:"not null;index"` // 외래키
	Title     string     `gorm:"type:text"`
	Content   string     `gorm:"type:text"`
	Status    string     `gorm:"type:varchar(20);not null;default:pending;index"`
	SentAt    *time.Time //nil 처리를 위해 포인터
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // 연관관계
}
