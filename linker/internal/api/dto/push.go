package dto

type PostPushReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type PostPushDTO struct {
	UserId  uint
	Title   string
	Content string
}

type CreateMessageDTO struct {
	UserId   uint
	Title    string
	Content  string
	SnsMsgId string
}
