package recommendation

import (
	"SingSong-Backend/internal/pkg"
	"SingSong-Backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomeRecommendation godoc
// @Summary      노래 추천 by 태그
// @Description  태그에 해당하는 노래를 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      HomeRequest  true  "태그 목록"
// @Success      200 {object} BaseResponse{data=[]HomeResponse} "성공"
// @Router       /recommend/tags [post]
func (pineconeHandler *PineconeHandler) HomeRecommendation(c *gin.Context) {
	// HomeRequest 형식으로 입력을 받습니다
	request := &usecase.HomeRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}

	homeResponses, err := pineconeHandler.recommendationUC.HomeRecommendation(c, request)
	if err != nil {
		pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
		return
	}

	pkg.BaseResponse(c, http.StatusOK, "ok", homeResponses)
	return
}
