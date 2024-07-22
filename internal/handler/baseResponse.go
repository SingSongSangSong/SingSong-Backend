package handler

// baseResponse 구조체 정의
type BaseResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 헬퍼 함수 작성
func NewBaseResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Message: message,
		Data:    data,
	}
}
