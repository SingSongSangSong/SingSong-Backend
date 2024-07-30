package handler

import (
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type songRecommendRequest struct {
	Songs []int `json:"songNumbers"`
}

type songRecommendResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

// RecommendBySongs godoc
// @Summary      노래 추천 by 노래 번호 목록
// @Description  노래 번호 목록을 보내면 유사한 노래들을 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      songRecommendRequest  true  "노래 번호 목록"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]songRecommendResponse} "성공"
// @Router       /recommend/songs [post]
func SongRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &songRecommendRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}
		// 리퀘스트에서 노래 목록을 슬라이스 형식으로 변환합니다.
		songs := make([]string, 0, len(request.Songs))
		for _, song := range request.Songs {
			songs = append(songs, strconv.Itoa(song))
		}

		// vectorid로 vector 조회하기
		res, err := idxConnection.FetchVectors(c, songs)
		if err != nil {
			log.Printf("Failed to fetch vectors, error: %+v", err)
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		returnSongs := make([]songRecommendResponse, 0)

		for i := 0; i < len(songs); i++ {
			wg.Add(1)
			go func(songID string) {
				defer wg.Done()
				vector, exists := res.Vectors[songID]

				if !exists {
					log.Printf("Vector with ID %s not found in response", songID)
					return
				}

				queryVector := vector.Values

				values, err := idxConnection.QueryByVectorValues(c, &pinecone.QueryByVectorValuesRequest{
					Vector:          queryVector,
					TopK:            uint32(20 / len(songs)),
					Filter:          nil,
					SparseValues:    nil,
					IncludeValues:   true,
					IncludeMetadata: true,
				})
				if err != nil {
					log.Printf("Failed to query by vector values, error: %+v", err)
					return
				}

				for j := 0; j < len(values.Matches); j++ {
					v := values.Matches[j].Vector
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

					mu.Lock()
					returnSongs = append(returnSongs, songRecommendResponse{
						songNumber,
						v.Metadata.Fields["song_name"].GetStringValue(),
						v.Metadata.Fields["singer_name"].GetStringValue(),
						koreanTags,
					})
					mu.Unlock()
				}
			}(songs[i])
		}

		wg.Wait()

		pkg.BaseResponse(c, http.StatusOK, "ok", returnSongs)
	}
}
