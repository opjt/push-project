package dto

type UpdateStatusReq struct {
	MsgId  uint64
	Status string
}

type UpdateStatusRes struct {
	MsgId uint64
}
