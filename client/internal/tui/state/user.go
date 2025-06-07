package state

type User struct {
	UserId string
	// SessionId string
}

func NewUser() *User {
	return &User{}
}
