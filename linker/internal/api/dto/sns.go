package dto

type SnsBody struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserId uint   `json:"user_id"`
}
