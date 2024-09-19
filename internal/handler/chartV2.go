package handler

import (
	"SingSong-Server/internal/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"
)

type V2RedisChartResponse struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"song_info_id"`
	TotalScore    float32 `json:"total_score"`
	New           int     `json:"new"`
	RankingChange int     `json:"ranking_change"`
	ArtistName    string  `json:"artist_name"`
	SongName      string  `json:"song_name"`
	SongNumber    int     `json:"song_number"`
	IsMR          int     `json:"is_mr"`
	ISLive        int     `json:"is_live"`
	Album         string  `json:"album"`
	Gender        string  `json:"gender"`
	AgeGroup      string  `json:"age_group"`
}

// ChartResponse 카멜케이스 구조체
type V2ChartSong struct {
	Ranking       int     `json:"ranking"`
	SongInfoId    int     `json:"songId"`
	TotalScore    float32 `json:"totalScore"`
	IsNew         bool    `json:"isNew"`
	RankingChange int     `json:"rankingChange"`
	ArtistName    string  `json:"artistName"`
	SongName      string  `json:"songName"`
	SongNumber    int     `json:"songNumber"`
	IsMR          bool    `json:"isMr"`
	IsLive        bool    `json:"isLive"`
	Album         string  `json:"album"`
}

type V2ChartOfKey struct {
	ChartKey string        `json:"chartKey"`
	Songs    []V2ChartSong `json:"songs"`
}

func convertOldToNewV2(old []V2RedisChartResponse) []V2ChartSong {
	var newCharts []V2ChartSong
	for _, o := range old {
		newCharts = append(newCharts, V2ChartSong{
			Ranking:       o.Ranking,
			SongInfoId:    o.SongInfoId,
			TotalScore:    o.TotalScore,
			IsNew:         o.New == 1, // 1,0 -> true/false로 변환
			RankingChange: o.RankingChange,
			ArtistName:    o.ArtistName,
			SongName:      o.SongName,
			SongNumber:    o.SongNumber,
			IsMR:          o.IsMR == 1,   // 1, 0 -> true/false로 변환
			IsLive:        o.ISLive == 1, // 1, 0 -> true/false로 변환
			Album:         o.Album,
		})
	}
	return newCharts
}

type V2TotalChartResponse struct {
	Time     string         `json:"time"`
	Gender   string         `json:"gender"`
	AgeGroup string         `json:"ageGroup"`
	UserKey  string         `json:"userKey"`
	Charts   []V2ChartOfKey `json:"charts"`
}

var V2ChartKey = []string{
	"MALE_ALL",
	"MALE_10",
	"MALE_20",
	"MALE_30",
	"MALE_40+",
	"FEMALE_ALL",
	"FEMALE_10",
	"FEMALE_20",
	"FEMALE_30",
	"FEMALE_40+",
	"MIXED_ALL",
	"MIXED_10",
	"MIXED_20",
	"MIXED_30",
	"MIXED_40+",
}

// GetChartV2 godoc
// @Summary      인기차트 조회(Version2)
// @Description  인기차트 조회(Version2)
// @Tags         Chart
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]V2TotalChartResponse} "성공"
// @Failure      400 {object} pkg.BaseResponseStruct{data=nil} "실패"
// @Router       /v2/chart [get]
// @Security BearerAuth
func GetChartV2(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		//성별 조회
		gender, exists := c.Get("gender")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - gender not found", nil)
			return
		}

		//
		birthYear, exists := c.Get("birthYear")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - birthyear not found", nil)
			return
		}

		// 전체 차트 만들기
		location, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot load location", nil)
			return
		}
		// Set the timezone for the current process
		time.Local = location
		currentTime := time.Now()

		var wholeCharts []V2ChartOfKey
		for index := range V2ChartKey {
			chartKey := V2ChartKey[index]
			redisKeyFormat := currentTime.Format("2006-01-02-15") + "-Hot_Trend_" + chartKey
			chart, err := rdb.Get(c, redisKeyFormat).Result()
			if err != nil {
				log.Printf("Redis GET error for key %s: %v", redisKeyFormat, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot find chart", nil)
				return
			}

			var redisChart []V2RedisChartResponse
			err = json.Unmarshal([]byte(chart), &redisChart)
			if err != nil {
				log.Printf("JSON Unmarshal error: %v", err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			responseChart := convertOldToNewV2(redisChart)
			if responseChart == nil {
				responseChart = []V2ChartSong{}
			}
			wholeCharts = append(wholeCharts, V2ChartOfKey{ChartKey: chartKey, Songs: responseChart})
		}

		ageGroup, userKey := createUserKey(gender, birthYear)

		response := V2TotalChartResponse{
			Time:     currentTime.Format("2006-01-02-15"),
			Gender:   gender.(string),
			AgeGroup: ageGroup,
			UserKey:  userKey,
			Charts:   wholeCharts,
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

func createUserKey(gender any, birthYear any) (string, string) {
	// 유저의 성별과 연령대로 KEY 만들기
	genderStr, ok := gender.(string)
	if !ok || (genderStr != "MALE" && genderStr != "FEMALE") {
		genderStr = "MIXED"
	}

	// birthYear를 string에서 int로 변환, 변환 실패 시 연령대는 "ALL"
	birthYearStr, ok := birthYear.(string)
	if !ok {
		// 변환 실패 시 기본값으로 처리
		return "ALL", fmt.Sprintf("%s_ALL", genderStr)
	}

	birthYearInt, err := strconv.Atoi(birthYearStr)
	if err != nil {
		// 변환 실패 시 기본값으로 처리
		return "ALL", fmt.Sprintf("%s_ALL", genderStr)
	}

	if birthYearInt == 0 {
		return "ALL", fmt.Sprintf("%s_ALL", genderStr)
	}

	currentYear := time.Now().Year()
	age := currentYear - birthYearInt + 1

	// 연령대 분류
	var ageGroup string
	switch {
	case age >= 10 && age < 20:
		ageGroup = "10"
	case age >= 20 && age < 30:
		ageGroup = "20"
	case age >= 30 && age < 40:
		ageGroup = "30"
	case age >= 40:
		ageGroup = "40+"
	default:
		ageGroup = "ALL"
	}

	// UserKey 생성 (예시: 성별_연령대)
	userKey := fmt.Sprintf("%s_%s", genderStr, ageGroup)
	return ageGroup, userKey
}
