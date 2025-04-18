package handler

import (
	"SingSong-Server/internal/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

// OldChartResponse 기존 구조체 (스네이크 케이스)
type OldChartResponse struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"song_info_id"`
	TotalScore    float32 `json:"total_score"`
	New           string  `json:"new"`
	RankingChange int     `json:"rankingChange"`
	ArtistName    string  `json:"artist_name"`
	SongName      string  `json:"song_name"`
	SongNumber    int     `json:"song_number"`
	IsMR          int     `json:"is_mr"`
}

// ChartResponse 카멜케이스 구조체
type ChartResponse struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"songId"`
	TotalScore    float32 `json:"totalScore"`
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

type TotalChartResponse struct {
	Time   string          `json:"time"`
	Gender string          `json:"gender"`
	Male   []ChartResponse `json:"male"`
	Female []ChartResponse `json:"female"`
}

// GetChart godoc
// @Summary      인기차트 조회
// @Description  인기차트 조회
// @Tags         Chart
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]TotalChartResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /v1/chart [get]
// @Security BearerAuth
func GetChart(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 성별 조회
		gender, exists := c.Get("gender")
		if !exists {
			pkg.SendToSentryWithStack(c, fmt.Errorf("gender not found in context"))
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not found", nil)
			return
		}

		var oldMaleCharts []OldChartResponse
		var oldFemaleCharts []OldChartResponse
		var maleCharts, femaleCharts []ChartResponse
		location, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot load location", nil)
			return
		}

		// Set the timezone for the current process
		time.Local = location
		currentTime := time.Now()

		// 남성 차트 조회
		maleFormattedTime := currentTime.Format("2006-01-02-15") + "-Hot_Trend_MALE"
		maleChart, err := rdb.Get(c, maleFormattedTime).Result()
		if errors.Is(err, redis.Nil) {
			log.Printf("No data found for male chart at %s", maleFormattedTime)
		} else if err != nil {
			log.Printf("Error retrieving male chart: %v", err)
		} else {
			if err := json.Unmarshal([]byte(maleChart), &oldMaleCharts); err != nil {
				log.Printf("Error parsing male chart JSON: %v", err)
			} else {
				maleCharts = convertOldToNew(oldMaleCharts)
			}
		}

		// 여성 차트 조회
		femaleFormattedTime := currentTime.Format("2006-01-02-15") + "-Hot_Trend_FEMALE"
		femaleChart, err := rdb.Get(c, femaleFormattedTime).Result()
		if errors.Is(err, redis.Nil) {
			log.Printf("No data found for female chart at %s", femaleFormattedTime)
		} else if err != nil {
			log.Printf("Error retrieving female chart: %v", err)
		} else {
			if err := json.Unmarshal([]byte(femaleChart), &oldFemaleCharts); err != nil {
				log.Printf("Error parsing female chart JSON: %v", err)
			} else {
				femaleCharts = convertOldToNew(oldFemaleCharts)
			}
		}

		totalChartResponse := TotalChartResponse{
			Time:   currentTime.Format("2006-01-02-15"),
			Gender: gender.(string),
			Male:   maleCharts,
			Female: femaleCharts,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", totalChartResponse)
	}
}
