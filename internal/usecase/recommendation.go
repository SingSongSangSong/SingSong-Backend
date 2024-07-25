package usecase

import (
	"SingSong-Backend/internal/app/user"
	"SingSong-Backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"strconv"
	"sync"
)

type RecommendationUseCase struct {
	redisModel    *model.RedisModel
	pineconeModel *model.PineconeModel
}

func NewRecommendationUseCase(redisModel *model.RedisModel, pineconeModel *model.PineconeModel) *RecommendationUseCase {
	return &RecommendationUseCase{redisModel: redisModel, pineconeModel: pineconeModel}
}

// Home 추천
type songHomeResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

type HomeRequest struct {
	Tags []string `json:"tags"`
}

type HomeResponse struct {
	Tag   string             `json:"tag"`
	Songs []songHomeResponse `json:"songs"`
}

func (recommendationUC *RecommendationUseCase) HomeRecommendation(c *gin.Context, request *HomeRequest) ([]HomeResponse, error) {
	// 한국어 태그가 들어오면 영어태그로 할당합니다
	englishTags, err := user.MapTagsKoreanToEnglish(request.Tags)
	log.Printf("englishTags: %v", englishTags)
	log.Printf("request.Tags: %v", request.Tags)
	if err != nil {
		//pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return nil, err
	}

	var homeResponses []HomeResponse
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

			// 노래들을 입력을 받습니다
			values, err := recommendationUC.pineconeModel.QueryPineconeWithTag(filterStruct)
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
				log.Printf("ssssArray: %v", ssssArray)
				koreanTags, err := user.MapTagsEnglishToKorean(ssssArray)

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

			koreanTag, err := user.MapTagEnglishToKorean(tag)
			mu.Lock()
			homeResponses = append(homeResponses, HomeResponse{
				Tag:   koreanTag,
				Songs: returnSongs,
			})
			mu.Unlock()
		}(i, tag)
	}
	wg.Wait()

	if overallErr != nil {
		return nil, overallErr
	}

	return homeResponses, nil
}

// song 추천
type SongRecommendRequest struct {
	Songs []int `json:"songs"`
}

type SongRecommendResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

func (recommendationUC *RecommendationUseCase) RecommendBySongs(c *gin.Context, request *SongRecommendRequest) ([]SongRecommendResponse, error) {
	// 리퀘스트에서 노래 목록을 슬라이스 형식으로 변환합니다.
	songs := make([]string, 0, len(request.Songs))
	for _, song := range request.Songs {
		songs = append(songs, strconv.Itoa(song))
	}

	// vectorid로 vector 조회하기
	res, err := recommendationUC.pineconeModel.FetchVectors(c, songs)
	if err != nil {
		log.Printf("Failed to fetch vectors, error: %+v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	returnSongs := make([]SongRecommendResponse, 0)

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

			values, err := recommendationUC.pineconeModel.QueryByVectorValues(c, &pinecone.QueryByVectorValuesRequest{
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
				koreanTags, err := user.MapTagsEnglishToKorean(ssssArray)

				if err != nil {
					log.Printf("Failed to convert tags to korean, error: %+v", err)
					koreanTags = []string{}
				}

				mu.Lock()
				returnSongs = append(returnSongs, SongRecommendResponse{
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
	return returnSongs, nil
}

//func (uc *RecommendationUseCase) Refresh(c *gin.Context, email string, koreanTag string) {
//	englishTag, err := user.MapTagKoreanToEnglish(koreanTag)
//
//	history := uc.getRefreshHistory(c, email, englishTag)
//
//	uc.setRefreshHistory(c, email, englishTag, history)
//}
//
//// Redis에서 새로고침 기록 저장하기 함수
//func (uc *RecommendationUseCase) setRefreshHistory(c *gin.Context, email string, englishTag string, songNumbers []int) {
//	//username:englishTag 형식으로 키 생성
//	//key 띄어써도 되나?
//	key := email + ":" + englishTag
//	response := uc.redisModel.Get(c, key)
//	if err := response.Err(); err != nil {
//		// case1: 홈 화면 요청인 경우
//		uc.redisModel.Set(c, key, songNumbers, 24*time.Hour) // 하루 지나면 만료
//		return
//	}
//
//	//case2: home 화면에서 새로고침 화면으로 진입 or 계속 새로고침 하는 경우
//	// response의 value 리스트에 songNumbers 추가
//	var existingNumbers []int
//	if err := response.Scan(&existingNumbers); err != nil {
//		// Scan 실패 시 로깅
//		log.Printf("기존 값 파싱 실패: %v", err)
//		return
//	}
//	updatedNumbers := append(existingNumbers, songNumbers...)
//	uc.redisModel.Set(c, key, updatedNumbers, 24*time.Hour)
//}
//
//// Redis에서 새로고침 기록 가져오기 함수
//func (uc *RecommendationUseCase) getRefreshHistory(c *gin.Context, email string, englishTag string) []int {
//	//username:englishTag 형식으로 키 생성
//	//key 띄어써도 되나?
//	key := email + ":" + englishTag
//	response := uc.redisModel.Get(c, key)
//	if err := response.Err(); err != nil {
//		return []int{}
//	}
//
//	var existingNumbers []int
//	if err := response.Scan(&existingNumbers); err != nil {
//		// Scan 실패 시 로깅
//		log.Printf("기존 값 파싱 실패: %v", err)
//		// 리턴?
//	}
//	return existingNumbers
//}
