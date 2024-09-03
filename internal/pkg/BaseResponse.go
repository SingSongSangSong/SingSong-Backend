package pkg

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// BaseResponse 구조체 정의
type BaseResponseStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 헬퍼 함수 작성
func BaseResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	span, _ := tracer.StartSpanFromContext(c.Request.Context(), "http.response")
	defer span.Finish()

	response := BaseResponseStruct{
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)

	span.SetTag("http.response.status", statusCode)
	span.SetTag("http.response.message", message)
}
