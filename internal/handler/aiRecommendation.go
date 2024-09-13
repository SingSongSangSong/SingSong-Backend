package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	pb "SingSong-Server/proto/userProfileRecommend"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
)

// songHomeResponse와 songResponse가 동일한 것으로 가정하고 사용
type songResponse struct {
	SongNumber   int    `json:"songNumber"`
	SongName     string `json:"songName"`
	SingerName   string `json:"singerName"`
	SongInfoId   int64  `json:"songId"`
	Album        string `json:"album"`
	IsMr         bool   `json:"isMr"`
	IsKeep       bool   `json:"isKeep"`
	KeepCount    int    `json:"keepCount"`
	CommentCount int    `json:"commentCount"`
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
func GetRecommendation(db *sql.DB) gin.HandlerFunc {
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
			log.Println("Gender not found in context - defaulting")
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

		// SongInfoId 리스트를 담을 빈 리스트 생성
		var songInfoIds []int64

		// gRPC response에서 SongInfoId만 추출
		for _, item := range response.SimilarItems {
			songInfoIds = append(songInfoIds, item.SongInfoId)
		}

		// []int64를 []interface{}로 변환
		songInfoInterface := make([]interface{}, len(songInfoIds))
		for i, v := range songInfoIds {
			songInfoInterface[i] = v
		}

		// Keep 여부 가져오기
		keepSongs, err := mysql.KeepSongs(qm.WhereIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 댓글 수 가져오기
		commentsCounts, err := mysql.Comments(qm.WhereIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Keep 수 가져오기
		keepCounts, err := mysql.KeepSongs(qm.WhereIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Keep 여부를 저장하는 맵 생성
		keepSongsMap := make(map[int64]bool)
		for _, keep := range keepSongs {
			keepSongsMap[keep.SongInfoID] = true // 해당 SongInfoId에 대해 Keep 여부를 기록
		}

		// 댓글 수를 저장하는 맵 생성
		commentsCountsMap := make(map[int64]int)
		for _, comment := range commentsCounts {
			commentsCountsMap[comment.SongInfoID]++
		}

		// Keep 수를 저장하는 맵 생성
		keepCountsMap := make(map[int64]int)
		for _, keep := range keepCounts {
			keepCountsMap[keep.SongInfoID]++
		}

		// gRPC response에서 가져온 SongInfoId를 기반으로 songInfoMap, keepSongsMap, commentsCountsMap, keepCountsMap을 활용
		for _, item := range response.SimilarItems {
			// 기본값으로 초기화
			isKeep := false
			commentCount := 0
			keepCount := 0

			if v, exists := keepSongsMap[item.SongInfoId]; exists {
				isKeep = v
			}
			if v, exists := commentsCountsMap[item.SongInfoId]; exists {
				commentCount = v
			}
			if v, exists := keepCountsMap[item.SongInfoId]; exists {
				keepCount = v
			}

			// userProfileRes.Songs에 추가
			userProfileRes.Songs = append(userProfileRes.Songs, songResponse{
				SongNumber:   int(item.SongNumber),
				SongName:     item.SongName,
				SingerName:   item.SingerName,
				SongInfoId:   item.SongInfoId,
				Album:        item.Album,
				IsMr:         item.IsMr,
				IsKeep:       isKeep,       // Keep 여부 추가
				CommentCount: commentCount, // 댓글 수 추가
				KeepCount:    keepCount,    // Keep 수 추가
			})
		}

		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", userProfileRes)
	}
}
