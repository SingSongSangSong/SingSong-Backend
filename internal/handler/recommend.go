package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"log"
	"net/http"
)

type RecommendRequest struct {
	Songs []string `json:"songs"`
}

type RecommendResponse struct {
	Songs []string `json:"songs"`
}

// GetRecommendation godoc
// @Summary      노래 추천 by 노래 번호 목록
// @Description  노래 번호 목록을 보내면 유사한 노래들을 추천합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      RecommendRequest  true  "노래 번호 목록"
// @Success      200 {object} RecommendResponse "성공"
// @Router       /recommend [post]
func (pineconeHandler *PineconeHandler) RegisterRecommendation(c *gin.Context) {
	request := &RecommendRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 리퀘스트에서 노래 목록을 슬라이스 형식으로 변환합니다.
	songs := make([]string, 0, len(request.Songs))
	for _, song := range request.Songs {
		songs = append(songs, song)
	}

	// vectorid로 vector 조회하기
	res, err := pineconeHandler.pinecone.FetchVectors(c, songs)
	if err != nil {
		log.Fatalf("Failed to fetch vectors, error: %+v", err)
	}

	returnSongs := make([]string, 0, len(res.Vectors))

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Iterate through the matches in the QueryVectorsResponse
		for j := 0; j < len(values.Matches); j++ {
			returnSongs = append(returnSongs, values.Matches[j].Vector.Id)
		}
	}

	// Returning the result as a JSON response
	c.JSON(http.StatusOK, RecommendResponse{returnSongs})
}

func (pineconeHandler *PineconeHandler) GetPineconeIndex(c *gin.Context) {
	idxs, err := pineconeHandler.pinecone.DescribeIndexStats(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"data":    idxs,
	})
}
