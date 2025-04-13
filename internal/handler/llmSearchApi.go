package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"math/rand"
	"net/http"
	"time"
)

type SearchResultForLLMResponse struct {
	SearchTexts []string `json:"searchTexts"`
}

// GetSearchResultsForLLM godoc
// @Summary      Get 3 Recent Search Results for LLM
// @Description  Get Recent 10 Search Results and provide Random 3 Search Texts for LLM
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=SearchResultForLLMResponse} "Success"
// @Router       /v1/recommend/recommendation/searchLog [get]
// @Security BearerAuth
func GetSearchResultsForLLM(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		llmSearch, err := mysql.LLMSearchLogs(
			qm.OrderBy("created_at DESC"),
			qm.Limit(10),
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
		randomTexts := getRandomElements(searchTexts, 3)

		pkg.BaseResponse(c, http.StatusOK, "success", SearchResultForLLMResponse{SearchTexts: randomTexts})
	}
}

// 랜덤으로 n개의 요소를 배열에서 선택
func getRandomElements(arr []string, n int) []string {
	// 시드 설정 (매번 다른 랜덤 값을 생성하기 위해)
	rand.Seed(time.Now().UnixNano())

	// 배열 길이가 n보다 작을 경우, 배열 자체를 반환
	if len(arr) < n {
		return arr
	}

	// n개의 랜덤 인덱스를 선택
	perm := rand.Perm(len(arr))[:n]

	// 선택된 인덱스를 바탕으로 배열에서 요소 추출
	var result []string
	for _, idx := range perm {
		result = append(result, arr[idx])
	}

	return result
}
