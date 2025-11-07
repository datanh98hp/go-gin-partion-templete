package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string // type alias string
// define error code
const (
	ErrorCodeBadRequest  ErrorCode = "BAD_REQUEST"
	ErrorCodeNotFound    ErrorCode = "NOT_FOUND"
	ErrorCodeInternal    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeConflict    ErrorCode = "CONFLICT"
	ErrCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

type AppError struct {
	Code    ErrorCode
	Erorr   error
	Message string
}
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message" binding:"omitempty"`
	Data    any    `json:"data" binding:"omitempty"`
}

func (err *AppError) Error() string {
	return err.Message
}
func NewError(message string, code ErrorCode) error {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewWrapError(message string, code ErrorCode, err error) error {
	return &AppError{
		Code:    code,
		Message: message,
		Erorr:   err,
	}
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok { // chuyển từ eror => appErr
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": CapitalLizeFirtCharacter(appErr.Message), //appErr.Message,
			"code":  appErr.Code,
		}
		if appErr.Erorr != nil {
			response["error_detail"] = appErr.Erorr.Error()
		}
		ctx.JSON(status, response)
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrorCodeInternal,
	})
}
func ResponseSuccess(ctx *gin.Context, status int, message string, data ...any) {
	resp := APIResponse{
		Status:  "success",
		Message: CapitalLizeFirtCharacter(message),
	}
	if len(data) > 0 && data[0] != nil {
		resp.Data = data[0]

	}
	ctx.JSON(status, resp)
}
func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}
func ResponseValidator(ctx *gin.Context, data any) {

	ctx.JSON(http.StatusBadRequest, data)
}
func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrorCodeBadRequest:
		return http.StatusBadRequest
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeInternal:
		return http.StatusInternalServerError
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
