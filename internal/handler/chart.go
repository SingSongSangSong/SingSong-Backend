package handler

import (
	"SingSong-Server/internal/pkg"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sync"
	"time"
)

// 기존 구조체 (스네이크 케이스)
type OldChartResponse struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"song_info_id"`
	TotalScore    float32 `json:"total_score"`
	Gender        string  `json:"gender"`
	BirthYear     int     `json:"birthYear"`
	New           string  `json:"new"`
	RankingChange int     `json:"rankingChange"`
	ArtistName    string  `json:"artist_name"`
	SongName      string  `json:"song_name"`
	SongNumber    int     `json:"song_number"`
	IsMR          int     `json:"is_mr"`
}

// 카멜케이스 구조체
type ChartResponse struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"songId"`
	TotalScore    float32 `json:"totalScore"`
	Gender        string  `json:"gender"`
	BirthYear     int     `json:"birthYear"`
	New           string  `json:"new"`
	RankingChange int     `json:"rankingChange"`
	ArtistName    string  `json:"artistName"`
	SongName      string  `json:"songName"`
	SongNumber    int     `json:"songNumber"`
	IsMR          int     `json:"isMr"`
}

func convertOldToNew(old []OldChartResponse) []ChartResponse {
	var newCharts []ChartResponse
	for _, o := range old {
		newCharts = append(newCharts, ChartResponse{
			Ranking:       o.Ranking,
			SongInfoId:    o.SongInfoId,
			TotalScore:    o.TotalScore,
			Gender:        o.Gender,
			BirthYear:     o.BirthYear,
			New:           o.New,
			RankingChange: o.RankingChange,
			ArtistName:    o.ArtistName,
			SongName:      o.SongName,
			SongNumber:    o.SongNumber,
			IsMR:          o.IsMR,
		})
	}
	return newCharts
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
		gender, exists := c.Get("gender")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - gender not found", nil)
			return
		}

		var oldMaleCharts []OldChartResponse
		var oldFemaleCharts []OldChartResponse
		var maleCharts, femaleCharts []ChartResponse
		currentTime := time.Now()

		var wg sync.WaitGroup
		wg.Add(2) // 두 개의 goroutine을 기다리기 위해 WaitGroup에 2를 추가

		go func() {
			defer wg.Done() // goroutine이 끝날 때 WaitGroup에 Done을 호출
			// 남성 차트 조회
			maleFormattedTime := currentTime.Format("2006-01-02-15") + "-Hot_Trend_MALE"
			maleChart, err := rdb.Get(c, maleFormattedTime).Result()
			if err != nil && err != redis.Nil {
				log.Printf("error - failed to get male chart: %v", err)
				return
			}
			// JSON 파싱
			if err := json.Unmarshal([]byte(maleChart), &oldMaleCharts); err != nil {
				log.Printf("Error parsing male chart JSON: %v", err)
			}
			maleCharts = convertOldToNew(oldMaleCharts)
		}()

		go func() {
			defer wg.Done() // goroutine이 끝날 때 WaitGroup에 Done을 호출
			// 여성 차트 조회
			femaleFormattedTime := currentTime.Format("2006-01-02-15") + "-Hot_Trend_FEMALE"
			femaleChart, err := rdb.Get(c, femaleFormattedTime).Result()
			if err != nil && err != redis.Nil {
				log.Printf("error - failed to get female chart: %v", err)
				return
			}
			// JSON 파싱
			if err := json.Unmarshal([]byte(femaleChart), &oldFemaleCharts); err != nil {
				log.Printf("Error parsing female chart JSON: %v", err)
			}
			femaleCharts = convertOldToNew(oldMaleCharts)
		}()

		wg.Wait() // 모든 goroutine이 끝날 때까지 대기

		// 결과 조합
		totalChart := map[string]interface{}{
			"Time":   currentTime.Format("2006-01-02-15"),
			"Gender": gender,
			"Male":   maleCharts,
			"Female": femaleCharts,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", totalChart)
	}
}
