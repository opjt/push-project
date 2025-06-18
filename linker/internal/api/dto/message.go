package dto

type CreateMessageDTO struct {
	UserId  uint64
	Title   string
	Content string
}

type UpdateMessageDTO struct {
	Id       uint64
	Status   string
	SnsMsgId string
}

type UpdateMessagesDTO struct {
	Ids    []uint64
	Column UpdateMessageField
}
type UpdateMessageField struct {
	Status   string
	SnsMsgId string
}
