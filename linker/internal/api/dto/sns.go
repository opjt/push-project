package dto

type SnsBody struct {
	MsgId  uint64 `json:"msg_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId uint64 `json:"user_id"`
}
