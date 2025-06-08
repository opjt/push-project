package dto

type CommonResponse[T any] struct {
	Error string `json:"error"`
	Data  T      `json:"data"`
}
