package state

type User struct {
	UserId    uint64
	Username  string
	SessionId string
}

func NewUser() *User {
	return &User{}
}
