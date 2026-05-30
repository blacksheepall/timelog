package router

import "net/http"

// ApiResponse is the frozen HTTP response envelope. Proto describes only Data.
type ApiResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
	Status  int         `json:"status"`
}

func SuccessResponse(data interface{}, message ...string) ApiResponse {
	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}
	return ApiResponse{
		Data:    data,
		Message: msg,
		Status:  http.StatusOK,
	}
}

func ErrorResponse(status int, message string) ApiResponse {
	return ApiResponse{
		Data:    nil,
		Message: message,
		Status:  status,
	}
}
