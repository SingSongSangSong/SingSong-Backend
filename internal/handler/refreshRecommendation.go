package handler

import (
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"math/rand"
	"net/http"
	"strconv"
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
	return func(c *gin.Context) {
		//todo: 유저 정보 필요 -> 어떤식으로 넘어오지?
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

		// todo: 레디스에서 조회해 온 숫자
		historySize := 20

		//refreshedSongs := make([]refreshResponse, 0, pageSize)
		vectorQuerySize := pageSize + historySize
		querySongs := make([]refreshResponse, 0, vectorQuerySize)

		dummyVector := make([]float32, 30)
		for i := range dummyVector {
			dummyVector[i] = rand.Float32()
		}
		log.Printf("querySize: ", vectorQuerySize)
		values, err := idxConnection.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
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

		// todo: 이전에 조회한 노래 빼고 pageSize개 반환

		pkg.BaseResponse(c, http.StatusOK, "ok", querySongs)
	}
}
