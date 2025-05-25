package dto

type SnsBody struct {
	MsgId  uint   `json:"msg_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId uint   `json:"user_id"`
}
