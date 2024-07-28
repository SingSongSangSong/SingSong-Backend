package handler

import (
	"SingSong-Server/internal/pkg"
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
// @Param        songs   body      refreshRequest  true  "태그"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]refreshResponse} "성공"
// @Router       /recommend/refresh [post]
func RefreshRecommendation(redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		historySongs := getRefreshHistory(c, redisClient, email, provider, englishTag)

		vectorQuerySize := pageSize + len(historySongs)
		values, err := queryVectorByTag(c, englishTag, idxConnection, vectorQuerySize)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - failed to query", nil)
			return
		}
		querySongs := extractSongInfo(vectorQuerySize, values)

		// 이전에 조회한 노래 빼고 상위 pageSize개 선택
		refreshedSongs := getTopSongsWithoutHistory(historySongs, querySongs)

		// 무한 새로고침 - (페이지의 끝일 때/노래 개수가 애초에 PageSize수보다 작을때) 부족한 노래 수만큼 refreshedSongs를 마저 채운다
		if len(refreshedSongs) < pageSize {
			refreshedSongs = fillSongsAgain(refreshedSongs, querySongs)
			// 기록 비우기
			historySongs = []int{}
		}

		// history 갱신
		for _, song := range refreshedSongs {
			historySongs = append(historySongs, song.SongNumber)
		}
		setRefreshHistory(c, redisClient, email, provider, historySongs, englishTag)

		pkg.BaseResponse(c, http.StatusOK, "ok", refreshedSongs)
	}
}

func getRefreshHistory(c *gin.Context, redisClient *redis.Client, email string, provider string, englishTag string) []int {
	key := generateRefreshKey(email, provider, englishTag)

	val, err := redisClient.Get(c, key).Result()
	if err == redis.Nil {
		return []int{}
	} else if err != nil {
		log.Printf("Failed to get history from Redis: %v", err)
		return []int{}
	}

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

func queryVectorByTag(c *gin.Context, englishTag string, idxConnection *pinecone.IndexConnection, vectorQuerySize int) (*pinecone.QueryVectorsResponse, error) {
	dummyVector := make([]float32, 30)
	for i := range dummyVector {
		dummyVector[i] = rand.Float32()*2 - 1 // -1 ~ 1
	}

	filterStruct := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"ssss": structpb.NewStringValue(englishTag),
			"MR":   structpb.NewBoolValue(false),
		},
	}
	values, err := idxConnection.QueryByVectorValues(c, &pinecone.QueryByVectorValuesRequest{
		Vector:          dummyVector,
		TopK:            uint32(vectorQuerySize),
		Filter:          filterStruct,
		SparseValues:    nil,
		IncludeValues:   true,
		IncludeMetadata: true,
	})
	return values, err
}

func extractSongInfo(vectorQuerySize int, values *pinecone.QueryVectorsResponse) []refreshResponse {
	querySongs := make([]refreshResponse, 0, vectorQuerySize)
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
	return querySongs
}

func getTopSongsWithoutHistory(historySongs []int, querySongs []refreshResponse) []refreshResponse {
	// golang에는 set이 없기 때문에 map을 구현해서 key만 사용하도록 했다
	refreshedSongs := make([]refreshResponse, 0, pageSize)
	historySet := toSet(historySongs)
	for _, song := range querySongs {
		if len(refreshedSongs) >= pageSize {
			break
		}
		if _, exists := historySet[song.SongNumber]; !exists {
			refreshedSongs = append(refreshedSongs, song)
		}
	}
	return refreshedSongs
}

func fillSongsAgain(refreshedSongs []refreshResponse, querySongs []refreshResponse) []refreshResponse {
	refreshedSongNumbers := make([]int, 0, len(refreshedSongs))
	for _, song := range refreshedSongs {
		refreshedSongNumbers = append(refreshedSongNumbers, song.SongNumber)
	}
	refreshedSet := toSet(refreshedSongNumbers)

	for _, song := range querySongs {
		if len(refreshedSongs) >= pageSize {
			break
		}
		// refreshedSongs 에 없는 곡으로 넣는다
		if _, exists := refreshedSet[song.SongNumber]; !exists {
			refreshedSongs = append(refreshedSongs, song)
		}
	}
	return refreshedSongs
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
