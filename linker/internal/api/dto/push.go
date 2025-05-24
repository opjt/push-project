package dto

type CreateMessageReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type CreateMessageDTO struct {
	UserID  uint   `json:"userid" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

type BroadcastMessageDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}
