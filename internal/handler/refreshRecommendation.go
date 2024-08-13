package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
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
	IsKeep     bool     `json:"isKeep"`
	SongInfoId int64    `json:"songId"`
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
// @Security BearerAuth
func RefreshRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		memberId, ok := value.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not type int64", nil)
			return
		}

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

		historySongs := getRefreshHistory(c, redisClient, memberId, englishTag)

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

		//list
		one, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 모든 KeepSongs 가져오기
		keepSongs, err := mysql.KeepSongs(qm.Where("keep_list_id = ?", one.KeepListID), qm.And("deleted_at IS NULL")).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Map으로 KeepSongs를 구성하여 존재 여부를 빠르게 확인
		isKeepMap := make(map[int]bool)
		for _, keepSong := range keepSongs {
			isKeepMap[keepSong.SongNumber] = true
		}

		// refreshSongs에 isKeep 여부 추가
		for i, song := range refreshedSongs {
			refreshedSongs[i].IsKeep = isKeepMap[song.SongNumber]
		}

		// SongInfoId 가져오기
		songNumbers := make([]interface{}, 0, len(refreshedSongs))
		for _, song := range refreshedSongs {
			songNumbers = append(songNumbers, song.SongNumber)
		}

		all, err := mysql.SongInfos(qm.WhereIn("song_number IN ?", songNumbers...)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		songTempIdMap := make(map[int]int64)
		for _, song := range all {
			songTempIdMap[song.SongNumber] = song.SongInfoID
		}

		// refreshSongs에 songTempId 추가
		for i, song := range refreshedSongs {
			refreshedSongs[i].SongInfoId = songTempIdMap[song.SongNumber]
		}

		// history 갱신
		for _, song := range refreshedSongs {
			historySongs = append(historySongs, song.SongNumber)
		}
		setRefreshHistory(c, redisClient, memberId, historySongs, englishTag)

		pkg.BaseResponse(c, http.StatusOK, "ok", refreshedSongs)
	}
}

func getRefreshHistory(c *gin.Context, redisClient *redis.Client, memberId int64, englishTag string) []int {
	key := generateRefreshKey(memberId, englishTag)

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

func generateRefreshKey(memberId int64, englishTag string) string {
	return "refresh:" + strconv.FormatInt(memberId, 10) + ":" + englishTag
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
		IncludeValues:   false,
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

func setRefreshHistory(c *gin.Context, redisClient *redis.Client, memberId int64, history []int, englishTag string) {
	key := generateRefreshKey(memberId, englishTag)

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
