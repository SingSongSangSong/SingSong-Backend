package handler

import (
	"SingSong-Server/internal/pkg"
	pb "SingSong-Server/proto/userProfileRecommend"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

// songHomeResponse와 songResponse가 동일한 것으로 가정하고 사용
type songResponse struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
}

type userProfileResponse struct {
	Songs []songResponse `json:"songs"`
}

// GetRecommendation godoc
// @Summary      AI가 골랐송
// @Description  사용자의 프로필을 기반으로 추천된 노래를 반환합니다. 페이지당 20개의 노래를 반환합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        pageId path int true "현재 조회할 노래 목록의 쪽수"
// @Success      200 {object} pkg.BaseResponseStruct{data=userProfileResponse} "성공"
// @Router       /v1/recommend/recommendation/{pageId} [get]
// @Security BearerAuth
func GetRecommendation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract pageId from the path
		page := c.Param("pageId")
		if page == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find pageId in path variable", nil)
			return
		}

		// Convert pageId to int32
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid pageId format", nil)
			return
		}

		// Get memberId from the middleware (assumed that the middleware sets the memberId)
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		gender, exists := c.Get("gender")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not found", nil)
			return
		}

		// Ensure memberId is cast to int64
		memberIdInt, ok := memberId.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - invalid memberId type", nil)
			return
		}

		// gRPC 서버에 연결
		conn, err := grpc.Dial("python-grpc:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewUserProfileClient(conn)

		// gRPC 요청 생성
		rpcRequest := &pb.ProfileRequest{
			MemberId: memberIdInt,
			Page:     int32(pageInt),
			Gender:   gender.(string),
		}

		// gRPC 요청 보내기
		response, err := client.CreateUserProfile(context.Background(), rpcRequest)
		if err != nil {
			log.Printf("Error calling gRPC: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Populate the userProfileResponse with gRPC response data
		userProfileRes := userProfileResponse{
			Songs: []songResponse{},
		}

		// Loop through the gRPC response to populate songResponse
		for _, item := range response.SimilarItems {
			userProfileRes.Songs = append(userProfileRes.Songs, songResponse{
				SongNumber: int(item.SongNumber),
				SongName:   item.SongName,
				SingerName: item.SingerName,
				SongInfoId: item.SongInfoId,
				Album:      item.Album,
				IsMr:       item.IsMr,
			})
		}

		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", userProfileRes)
	}
}
