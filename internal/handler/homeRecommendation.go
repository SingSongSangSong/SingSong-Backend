package handler

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type songHomeResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

// HomeRequest는 추천 요청 구조체입니다.
type HomeRequest struct {
	Tags []string `json:"tags"`
}

type HomeResponse struct {
	Tag   string             `json:"tag"`
	Songs []songHomeResponse `json:"songs"`
}

// HomeRecommendation godoc
// @Summary      노래 추천 by 태그
// @Description  태그에 해당하는 노래를 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      HomeRequest  true  "태그 목록"
// @Success      200 {object} BaseResponse{data=[]HomeResponse} "성공"
// @Router       /recommend/tags [post]
func (pineconeHandler *PineconeHandler) HomeRecommendation(c *gin.Context) {
	// HomeRequest 형식으로 입력을 받습니다
	request := &HomeRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}

	// 한국어 태그가 들어오면 영어태그로 할당합니다
	englishTags, err := mapTagsKoreanToEnglish(request.Tags)
	if err != nil {
		BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}
	var homeResponses []HomeResponse
	var wg sync.WaitGroup
	var mu sync.Mutex
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
				},
			}
			// 입력받을 노래들의 리스트를 할당합니다
			returnSongs := make([]songHomeResponse, 0, len(englishTags))

			// 노래들을 입력을 받습니다
			values, err := pineconeHandler.queryPineconeWithTag(filterStruct)
			if err != nil {
				BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// 받아온 입력들의 아이디 및 다른 값들을 할당합니다
			for _, match := range values.Matches {
				v := match.Vector
				songNumber, err := strconv.Atoi(v.Id)
				if err != nil {
					log.Printf("Failed to convert ID to int, error: %+v", err)
				}
				koreanTags, err := mapTagsEnglishToKorean(parseTags(v.Metadata.Fields["ssss"].GetStringValue()))
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

			koreanTag, err := mapTagEnglishToKorean(tag)
			mu.Lock()
			homeResponses = append(homeResponses, HomeResponse{
				Tag:   koreanTag,
				Songs: returnSongs,
			})
			mu.Unlock()
		}(i, tag)
	}
	wg.Wait()

	BaseResponse(c, http.StatusOK, "ok", homeResponses)
	return
}
