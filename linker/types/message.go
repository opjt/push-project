package types

const (
	StatusPending  = "pending"  // 메시지 발행 후, 아직 sender가 처리 전인 상태
	StatusSending  = "sending"  // sender가 메시지를 받아 실제 전송 작업 중인 상태
	StatusDeferred = "deferred" // 클라이언트가 세션이 연결되어 있지 않아 보류 중인 상태
	StatusSent     = "sent"     // 메시지 전송 완료 상태, client에서 확인된 경우.
	StatusFailed   = "failed"   // 메시지 전송 실패 상태
)
