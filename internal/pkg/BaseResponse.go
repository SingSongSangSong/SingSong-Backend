package pkg

import "github.com/gin-gonic/gin"

// BaseResponse 구조체 정의
type BaseResponseStruct struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 헬퍼 함수 작성
func NewBaseResponse(statusCode int, message string, data interface{}) BaseResponseStruct {
	return BaseResponseStruct{
		Message: message,
		Data:    data,
	}
}

// 헬퍼 함수 작성
func BaseResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := BaseResponseStruct{
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}
