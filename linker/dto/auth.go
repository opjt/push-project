package dto

// 클라이언트 로그인 요청시 req body.
type AuthLoginReq struct {
	UserId string `json:"userId" binding:"required"`
}

type AuthLoginRes struct {
	UserId    string
	SessionId string
}

type LoginResult struct {
	UserId    string
	SessionId string
}
