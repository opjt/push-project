package dto

// SQS 메세지 구조
type SqsEnvelope struct {
	Type              string `json:"Type"`
	MessageId         string `json:"MessageId"`
	TopicArn          string `json:"TopicArn"`
	Message           string `json:"Message"` // 여기에 내부 메시지 JSON이 감싸져 있음
	Timestamp         string `json:"Timestamp"`
	MessageAttributes map[string]struct {
		Type  string `json:"Type"`
		Value string `json:"Value"`
	} `json:"MessageAttributes"`
}

type PushMessage struct {
	MsgID  int    `json:"msg_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"user_id"`
}

type UpdateMessageStatus struct {
	Id       uint64
	Status   string
	SnsMsgId string
}
