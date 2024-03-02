package responses

import "net/http"

var (
	ErrAPIFailed      = new(http.StatusInternalServerError, 1000, "API failed", nil)
	ErrValidateFailed = new(http.StatusBadRequest, 1001, "validate failed", nil)
)

var (
	PostSuccess = new(http.StatusCreated, 0, "success", nil)
)

type Response struct {
	HTTPStatusCode int    `json:"-"`
	Code           int    `json:"code"`
	Message        string `json:"message"`
	Data           any    `json:"data"`
}

func new(httpStatusCode, code int, message string, data any) Response {
	return Response{
		HTTPStatusCode: httpStatusCode,
		Code:           code,
		Message:        message,
		Data:           data,
	}
}

func (r Response) Error() string {
	return r.Message
}
