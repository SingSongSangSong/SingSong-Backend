package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type ChartResponse struct {
	ChartNumber int    `json:"chartNumber"`
	ChartName   string `json:"chartName"`
}

// GetChart godoc
// @Summary      인기차트 조회
// @Description  인기차트 조회
// @Tags         Chart
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]ChartResponse} "성공"
// @Router       /chart [get]
// @Security BearerAuth
func GetChart(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 성별 조회
		gender, err := c.Get("gender")
		if err != true {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}
		// 생년 조회
		//birthYear, err := c.Get("birthYear")
		//if err != true {
		//	pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
		//	return
		//}

		// 인기차트 조회
		if gender == "MALE" {
			// 남성 인기차트 조회
			rdb.Get(c, "maleChart")
		}

	}
}
