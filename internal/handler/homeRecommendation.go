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
	"sync"
)

// Home 추천
type songHomeResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

type homeRequest struct {
	Tags []string `json:"tags"`
}

type homeResponse struct {
	Tag   string             `json:"tag"`
	Songs []songHomeResponse `json:"songs"`
}

// HomeRecommendation godoc
// @Summary      노래 추천 by 태그
// @Description  태그에 해당하는 노래를 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      homeRequest  true  "태그 목록"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]homeResponse} "성공"
// @Router       /recommend/home [post]
func HomeRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &homeRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// 한국어 태그가 들어오면 영어태그로 할당합니다
		englishTags, err := MapTagsKoreanToEnglish(request.Tags)
		log.Printf("englishTags: %v", englishTags)
		log.Printf("request.Tags: %v", request.Tags)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		var homeResponses []homeResponse
		var wg sync.WaitGroup
		var mu sync.Mutex
		var overallErr error

		// 각 태그에 대해서 돌면서 값을 가져온다!
		for i, tag := range englishTags {
			// 각 태그에 대해서 고루틴을 실행할때 WaitGroup을 추가하여 모두 마무리가 되었을때 넘어가도록 한다
			wg.Add(1)
			go func(i int, tag string) {
				defer wg.Done()

				// structpb.Struct 생성
				filterStruct := &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"ssss": structpb.NewStringValue(tag),
						"MR":   structpb.NewBoolValue(false),
					},
				}
				// 입력받을 노래들의 리스트를 할당합니다
				returnSongs := make([]songHomeResponse, 0, len(englishTags))

				// Define a dummy vector (e.g., zero vector) for the query
				dummyVector := make([]float32, 30) // Assuming the vector length is 1536, adjust as necessary
				for i := range dummyVector {
					dummyVector[i] = rand.Float32() //random vector
				}

				mu.Lock()

				// 쿼리 요청을 보냅니다.
				values, err := idxConnection.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
					Vector:          dummyVector,
					TopK:            20,
					Filter:          filterStruct,
					SparseValues:    nil,
					IncludeValues:   true,
					IncludeMetadata: true,
				})

				mu.Unlock()

				if err != nil {
					//pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
					// 에러 발생 시 전체 함수를 중단하지 않고 로그를 남기고 전체 에러 상태를 설정
					log.Printf("QueryPineconeWithTag error for tag %s: %+v", tag, err)
					mu.Lock()
					overallErr = err
					mu.Unlock()
					return
				}

				// 받아온 입력들의 아이디 및 다른 값들을 할당합니다
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
					returnSongs = append(returnSongs, songHomeResponse{
						SongNumber: songNumber,
						SongName:   v.Metadata.Fields["song_name"].GetStringValue(),
						SingerName: v.Metadata.Fields["singer_name"].GetStringValue(),
						Tags:       koreanTags,
					})
				}

				koreanTag, err := MapTagEnglishToKorean(tag)
				mu.Lock()
				homeResponses = append(homeResponses, homeResponse{
					Tag:   koreanTag,
					Songs: returnSongs,
				})
				mu.Unlock()
			}(i, tag)
		}
		wg.Wait()

		if overallErr != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", homeResponses)
	}
}
