package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

// GetSearchResultsForLLMV2 godoc
// @Summary      Get 10 Recent Search Results for LLM
// @Description  Get Recent 10 Search Results for LLM
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=SearchResultForLLMResponse} "Success"
// @Router       /v2/recommend/recommendation/searchLog [get]
// @Security BearerAuth
func GetSearchResultsForLLMV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		llmSearch, err := mysql.LLMSearchLogs(
			qm.OrderBy("created_at DESC"),
			qm.Limit(15),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		var searchTexts []string
		for _, l := range llmSearch {
			searchTexts = append(searchTexts, l.SearchText)
		}

		// Get random 3 search texts
		// 랜덤으로 3개의 텍스트 선택
		randomTexts := getRandomElements(searchTexts, 10)

		pkg.BaseResponse(c, http.StatusOK, "success", SearchResultForLLMResponse{SearchTexts: randomTexts})
	}
}
