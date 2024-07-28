package handler

import (
	"SingSong-Server/internal/pkg"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type refreshRequest struct {
	Tag string `json:"tag"`
}

type refreshResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

var (
	pageSize = 20
)

// RefreshRecommendation godoc
// @Summary      새로고침 노래 추천
// @Description  태그에 해당하는 노래를 새로고침합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      refreshRequest  true  "태그 목록"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]refreshResponse} "성공"
// @Router       /recommend/refresh [post]
func RefreshRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	f := func(c *gin.Context) {
		//todo: 유저 정보 필요 -> accesstoken에서 추출
		//일단 userEmail은 test@test.com 으로, provider는 kakao로 가정
		email := "test@test.com"
		provider := "kakao"

		request := &refreshRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		englishTag, err := MapTagKoreanToEnglish(request.Tag)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}
		log.Printf("englishTag: %v", englishTag)

		filterStruct := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"ssss": structpb.NewStringValue(englishTag),
				"MR":   structpb.NewBoolValue(false),
			},
		}

		historySongs := getRefreshHistory(c, redisClient, email, provider, englishTag)
		log.Printf("historySongs: %v", len(historySongs))
		vectorQuerySize := pageSize + len(historySongs)
		querySongs := make([]refreshResponse, 0, vectorQuerySize)
		dummyVector := make([]float32, 30)
		for i := range dummyVector {
			dummyVector[i] = rand.Float32()
		}
		log.Printf("querySize: ", vectorQuerySize)
		values, err := idxConnection.QueryByVectorValues(c, &pinecone.QueryByVectorValuesRequest{
			Vector:          dummyVector,
			TopK:            uint32(vectorQuerySize),
			Filter:          filterStruct,
			SparseValues:    nil,
			IncludeValues:   true,
			IncludeMetadata: true,
		})

		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - failed to query", nil)
			return
		}

		log.Printf("조회 벡터 크기: ", strconv.Itoa(len(values.Matches)))

		for _, match := range values.Matches {
			v := match.Vector
			songNumber, err := strconv.Atoi(v.Id)
			if err != nil {
				log.Printf("Failed to convert ID to int, error: %+v", err)
			}

			ssssField := v.Metadata.Fields["ssss"].GetListValue().AsSlice()
			ssssArray := make([]string, len(ssssField))
			for i, eTag := range ssssField {
				ssssArray[i] = eTag.(string)
			}
			koreanTags, err := MapTagsEnglishToKorean(ssssArray)

			if err != nil {
				log.Printf("Failed to convert tags to korean, error: %+v", err)
				koreanTags = []string{}
			}
			querySongs = append(querySongs, refreshResponse{
				SongNumber: songNumber,
				SongName:   v.Metadata.Fields["song_name"].GetStringValue(),
				SingerName: v.Metadata.Fields["singer_name"].GetStringValue(),
				Tags:       koreanTags,
			})
		}

		// 이전에 조회한 노래 빼고 상위 pageSize개 반환
		// golang에는 set이 없기 때문에 map을 구현해서 key만 사용하도록 했다
		historySet := toSet(historySongs)
		refreshedSongs := make([]refreshResponse, 0, pageSize)
		for _, song := range querySongs {
			if len(refreshedSongs) >= pageSize {
				break
			}
			if _, exists := historySet[song.SongNumber]; !exists {
				refreshedSongs = append(refreshedSongs, song)
			}
		}

		// todo: 비동기?
		// 기존 history + 이번에 새로고침된 곡들 덧붙여서 저장
		for _, song := range refreshedSongs {
			historySongs = append(historySongs, song.SongNumber)
		}
		setRefreshHistory(c, redisClient, email, provider, historySongs, englishTag)

		// todo: 이미 다 한번씩 조회했었다면? -> 다시 처음부터

		pkg.BaseResponse(c, http.StatusOK, "ok", refreshedSongs)
	}
	return f
}

func getRefreshHistory(c *gin.Context, redisClient *redis.Client, email string, provider string, englishTag string) []int {
	key := generateRefreshKey(email, provider, englishTag)

	val, err := redisClient.Get(c, key).Result()
	if err == redis.Nil {
		return []int{}
	} else if err != nil {
		// 다른 에러가 발생한 경우 로그를 남기고 빈 슬라이스를 반환합니다.
		log.Printf("Failed to get history from Redis: %v", err)
		return []int{}
	}

	// JSON 데이터를 슬라이스로 역직렬화합니다.
	var history []int
	err = json.Unmarshal([]byte(val), &history)
	if err != nil {
		log.Printf("Failed to unmarshal history: %v", err)
		return []int{}
	}

	return history
}

func generateRefreshKey(email string, provider string, englishTag string) string {
	return "refresh:" + email + ":" + provider + ":" + englishTag
}

func setRefreshHistory(c *gin.Context, redisClient *redis.Client, email string, provider string, history []int, englishTag string) {
	key := generateRefreshKey(email, provider, englishTag)

	historyJSON, err := json.Marshal(history)
	if err != nil {
		log.Printf("Failed to marshal history: %v", err)
		return
	}

	redisClient.Set(c, key, historyJSON, 30*time.Minute)
}

func toSet(slice []int) map[int]struct{} {
	set := make(map[int]struct{}, len(slice))
	for _, e := range slice {
		set[e] = struct{}{}
	}
	return set
}
