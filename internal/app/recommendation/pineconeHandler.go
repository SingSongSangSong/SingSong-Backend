package recommendation

import (
	"SingSong-Backend/internal/usecase"
)

type PineconeHandler struct {
	recommendationUC *usecase.RecommendationUseCase
}

// NewPineconeHandler는 PineconeHandler를 초기화하는 함수입니다.
func NewPineconeHandler(recommendationUC *usecase.RecommendationUseCase) (*PineconeHandler, error) {
	pcHandler := &PineconeHandler{recommendationUC: recommendationUC}
	return pcHandler, nil
}
