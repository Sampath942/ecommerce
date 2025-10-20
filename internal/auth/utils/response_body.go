package utils

type Response struct {
	ResponseMessage string `json:"message"`
	ResponseDetails any `json:"details"`
}