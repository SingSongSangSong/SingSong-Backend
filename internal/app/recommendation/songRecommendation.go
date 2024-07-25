package recommendation

import (
	"SingSong-Backend/internal/pkg"
	"SingSong-Backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RecommendBySongs godoc
// @Summary      노래 추천 by 노래 번호 목록
// @Description  노래 번호 목록을 보내면 유사한 노래들을 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      usecase.SongRecommendRequest  true  "노래 번호 목록"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]usecase.SongRecommendResponse} "성공"
// @Router       /recommend [post]
func (pineconeHandler *PineconeHandler) RecommendBySongs(c *gin.Context) {
	request := &usecase.SongRecommendRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}
	returnSongs, err := pineconeHandler.recommendationUC.RecommendBySongs(c, request)
	if err != nil {
		pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
		return
	}

	pkg.BaseResponse(c, http.StatusOK, "ok", returnSongs)
	return
}
