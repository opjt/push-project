package dto

type PostPushReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type PostPushDTO struct {
	UserId  uint64
	Title   string
	Content string
}
