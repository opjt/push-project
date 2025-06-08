package dto

// 클라이언트 로그인 요청시 req body.
type AuthLoginReq struct {
	Username string `json:"username" binding:"required"`
}

type AuthLoginRes struct {
	Username  string
	UserId    uint64
	SessionId string
}

type LoginResult struct {
	Username  string
	UserId    uint64
	SessionId string
}
