package dto

// 메세지 변경 요청을 위한 dto(http)
type UpdateStatusReq struct {
	MsgId  uint64
	Status string
}

// 메세지 변경 요청 결과
type UpdateStatusRes struct {
	MsgId uint64
}
