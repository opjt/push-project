package dto

type CreateMessageDTO struct {
	UserId  uint
	Title   string
	Content string
}

type UpdateMessageDTO struct {
	Id     uint
	Status string
}
