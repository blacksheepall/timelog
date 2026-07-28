package router

import "net/http"

// ApiResponse is the frozen HTTP response envelope for every /api response,
// mirrored by web/src/types/api.ts. Proto describes only Data.
//
// Contract (enforced by TestEnvelopeContract and the golden tests in
// contract_test.go):
//   - Success: HTTP 200, body {data: T, message: "...", status: 200}.
//     Command endpoints (delete/complete/...) use data: null.
//   - Error: non-2xx, body {data: null, message: "<human readable>", status: <same code>}.
//     The frontend surfaces message via err.response?.data?.message; the only
//     status it special-cases is 401 (logout + redirect to /login).
//   - JSON field names are snake_case everywhere on the wire (Go json tags on
//     generated proto structs + ts-proto snakeToCamel=false).
//
// Any change here must be mirrored in web/src/types/api.ts and the docs
// (docs/development-cn.md "API 请求/响应约定").
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
