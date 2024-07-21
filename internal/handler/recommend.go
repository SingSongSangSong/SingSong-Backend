package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RecommendRequest struct {
	Songs []int `json:"songs"`
}

type RecommendResponse struct {
	Songs []SongResponse `json:"songs"`
}

type SongResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
}

// GetRecommendation godoc
// @Summary      노래 추천 by 노래 번호 목록
// @Description  노래 번호 목록을 보내면 유사한 노래들을 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      RecommendRequest  true  "노래 번호 목록"
// @Success      200 {object} BaseResponse{data=RecommendResponse} "성공"
// @Router       /recommend [post]
func (pineconeHandler *PineconeHandler) GetSongRecommendation(c *gin.Context) {
	request := &RecommendRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse("error", nil))
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 리퀘스트에서 노래 목록을 슬라이스 형식으로 변환합니다.
	songs := make([]string, 0, len(request.Songs))
	for _, song := range request.Songs {
		songs = append(songs, strconv.Itoa(song))
	}

	// vectorid로 vector 조회하기
	res, err := pineconeHandler.pinecone.FetchVectors(c, songs)
	if err != nil {
		log.Fatalf("Failed to fetch vectors, error: %+v", err)
	}

	returnSongs := make([]SongResponse, 0, len(res.Vectors))

	for i := 0; i < len(songs); i++ {
		vector, exists := res.Vectors[songs[i]]

		if !exists {
			log.Printf("Vector with ID %s not found in response", songs[i])
			continue
		}

		queryVector := vector.Values

		values, err := pineconeHandler.pinecone.QueryByVectorValues(c, &pinecone.QueryByVectorValuesRequest{
			Vector:          queryVector,
			TopK:            uint32(20 / len(songs)),
			Filter:          nil,
			SparseValues:    nil,
			IncludeValues:   true,
			IncludeMetadata: true,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, NewBaseResponse("error", nil))
			//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Iterate through the matches in the QueryVectorsResponse
		for j := 0; j < len(values.Matches); j++ {
			v := values.Matches[j].Vector
			songNumber, err := strconv.Atoi(v.Id)
			if err != nil {
				log.Fatalf("Failed to convert ID to int, error: %+v", err)
			}

			returnSongs = append(returnSongs, SongResponse{
				songNumber,
				v.Metadata.Fields["song_name"].GetStringValue(),
				v.Metadata.Fields["singer_name"].GetStringValue(),
				parseTags(v.Metadata.Fields["ssss"].GetStringValue()),
			})
		}
	}

	// Returning the result as a JSON response
	c.JSON(http.StatusOK, NewBaseResponse("ok", RecommendResponse{returnSongs}))
}

func parseTags(tags string) []string {
	tagList := strings.Split(tags, ",")
	for i := range tagList {
		tagList[i] = strings.TrimSpace(tagList[i])
	}
	return tagList
}

func (pineconeHandler *PineconeHandler) GetPineconeIndex(c *gin.Context) {
	idxs, err := pineconeHandler.pinecone.DescribeIndexStats(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse("error", err.Error()))
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, NewBaseResponse("ok", idxs))
}

// HomeRequest는 추천 요청 구조체입니다.
type HomeRequest struct {
	Tags []string `json:"tags"`
}

// HomeRecommendation은 추천 요청을 처리하는 핸들러 함수입니다.
func (pineconeHandler *PineconeHandler) HomeRecommendation(c *gin.Context) {
	request := &HomeRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewBaseResponse("error", nil))
		//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 필터링 조건을 설정합니다.
	filterConditions := make([]*structpb.Value, len(request.Tags))
	for i, tag := range request.Tags {
		filterConditions[i] = structpb.NewStructValue(&structpb.Struct{
			Fields: map[string]*structpb.Value{
				"ssss": {
					Kind: &structpb.Value_StringValue{
						StringValue: tag,
					},
				},
			},
		})
	}

	filterStruct := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"$or": {
				Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{
						Values: filterConditions,
					},
				},
			},
		},
	}

	// Define a dummy vector (e.g., zero vector) for the query
	dummyVector := make([]float32, 30) // Assuming the vector length is 1536, adjust as necessary

	// 쿼리 요청을 보냅니다.
	values, err := pineconeHandler.pinecone.QueryByVectorValues(context.Background(), &pinecone.QueryByVectorValuesRequest{
		Vector:          dummyVector,
		TopK:            100,
		Filter:          filterStruct,
		SparseValues:    nil,
		IncludeValues:   true,
		IncludeMetadata: true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewBaseResponse("error", nil))
		//c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 반환된 벡터의 ID를 수집합니다.
	returnSongs := make([]string, 0, len(values.Matches))
	for _, match := range values.Matches {
		returnSongs = append(returnSongs, match.Vector.Id)
	}

	c.JSON(http.StatusOK, NewBaseResponse("ok", returnSongs))
}
