package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

// Home 추천
type songHomeResponse struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
	IsLive     bool   `json:"isLive"`
	MelonLink  string `json:"melonLink"`
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
// @Router       /v1/recommend/home [post]
func HomeRecommendation(db *sql.DB, redisClient *redis.Client, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &homeRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// 한국어 태그가 들어오면 영어태그로 할당합니다
		englishTags, err := MapTagsKoreanToEnglish(request.Tags)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// tagMap
		tagMap := make(map[string][]interface{})
		for _, tag := range englishTags {
			tagMap[tag] = nil
		}

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

				// Define a dummy vector (e.g., zero vector) for the query
				dummyVector := make([]float32, conf.VectorDBConfigInstance.DIMENSION) // Assuming the vector length is 1536, adjust as necessary
				for i := range dummyVector {
					dummyVector[i] = rand.Float32() //random vector
				}

				// 쿼리 요청을 보냅니다.
				values, err := idxConnection.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
					Vector:          dummyVector,
					TopK:            20,
					Filter:          filterStruct,
					SparseValues:    nil,
					IncludeValues:   false,
					IncludeMetadata: false,
				})

				if err != nil {
					//pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
					// 에러 발생 시 전체 함수를 중단하지 않고 로그를 남기고 전체 에러 상태를 설정
					log.Printf("QueryPineconeWithTag error for tag %s: %+v", tag, err)
					overallErr = err
					return
				}

				ids := make([]interface{}, 0, len(values.Matches))

				for _, match := range values.Matches {
					v := match.Vector
					songInfoId, err := strconv.Atoi(v.Id)
					if err != nil {
						log.Printf("Failed to convert ID to int, error: %+v", err)
					}
					ids = append(ids, int64(songInfoId))
				}

				mu.Lock()
				tagMap[tag] = ids
				mu.Unlock()
			}(i, tag)
		}
		wg.Wait()

		if overallErr != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		//tagMap의 value들을 합쳐 db에서 가져옵니다.
		ids := make([]interface{}, 0, 20*len(tagMap))
		for _, idsInTag := range tagMap {
			ids = append(ids, idsInTag...)
		}
		songInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", ids...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		homeResponses := make([]homeResponse, 0, len(tagMap))
		for tag, ids := range tagMap {
			korean, err := MapTagEnglishToKorean(tag)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			homeResponse := homeResponse{
				Tag:   korean,
				Songs: make([]songHomeResponse, 0, 20),
			}

			for _, songInfo := range songInfos {
				for _, id := range ids {
					if songInfo.SongInfoID == id.(int64) {
						homeResponse.Songs = append(homeResponse.Songs, songHomeResponse{
							SongNumber: songInfo.SongNumber,
							SongName:   songInfo.SongName,
							SingerName: songInfo.ArtistName,
							SongInfoId: songInfo.SongInfoID,
							Album:      songInfo.Album.String,
							IsMr:       songInfo.IsMR.Bool,
							IsLive:     songInfo.IsLive.Bool,
							MelonLink:  CreateMelonLinkByMelonSongId(songInfo.MelonSongID),
						})
					}
				}
			}
			homeResponses = append(homeResponses, homeResponse)
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", homeResponses)
	}
}
