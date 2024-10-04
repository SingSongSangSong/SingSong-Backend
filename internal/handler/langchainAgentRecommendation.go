package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	pb "SingSong-Server/proto/langchainAgentRecommend"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

// LangchainAgentRecommedation godoc
// @Summary      LLM으로 검색하기
// @Description  LLM의 사용자 입력을 토대로 추천된 노래를 반환합니다. 10개의 노래를 반환합니다
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        input   body      LlmRequest  true  "인풋"
// @Success      200 {object} pkg.BaseResponseStruct{data=userProfileResponse} "성공"
// @Router       /v1/recommend/recommendation/langchainAgent [post]
// @Security BearerAuth
func LangchainAgentRecommedation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get memberId from the middleware (assumed that the middleware sets the memberId)
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// Get the input from the request body
		llmRequest := LlmRequest{}
		if err := c.ShouldBindJSON(&llmRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// Ensure memberId is cast to int64
		memberIdInt, ok := memberId.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - invalid memberId type", nil)
			return
		}

		// gRPC 서버에 연결
		conn, err := grpc.Dial(GrpcAddr+":50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewLangchainAgentRecommendClient(conn)

		// gRPC 요청 생성
		rpcRequest := &pb.LangchainAgentRequest{
			MemberId: memberIdInt,
			Command:  llmRequest.UserInput,
		}

		// gRPC 요청 보내기
		response, err := client.GetLangchainAgentRecommendation(context.Background(), rpcRequest)
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
		for _, item := range response.SearchResult {
			songInfoIds = append(songInfoIds, item.SongInfoId)
		}

		// []int64를 []interface{}로 변환
		songInfoInterface := make([]interface{}, len(songInfoIds))
		for i, v := range songInfoIds {
			songInfoInterface[i] = v
		}

		// MelonSongId 가져오기
		songInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// MelonSongId를 저장하는 맵 생성
		songInfoMap := make(map[int64]*mysql.SongInfo)
		for _, songInfo := range songInfos {
			songInfoMap[songInfo.SongInfoID] = songInfo
		}

		// Loop through the gRPC response to populate songResponse
		for _, item := range response.SearchResult {
			userProfileRes.Songs = append(userProfileRes.Songs, songResponse{
				SongNumber: songInfoMap[item.SongInfoId].SongNumber,
				SongName:   songInfoMap[item.SongInfoId].SongName,
				SingerName: songInfoMap[item.SongInfoId].ArtistName,
				SongInfoId: songInfoMap[item.SongInfoId].SongInfoID,
				Album:      songInfoMap[item.SongInfoId].Album.String,
				IsMr:       songInfoMap[item.SongInfoId].IsMR.Bool,
				IsLive:     songInfoMap[item.SongInfoId].IsLive.Bool,
				MelonLink:  CreateMelonLinkByMelonSongId(songInfoMap[item.SongInfoId].MelonSongID),
			})
		}

		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", userProfileRes)
	}
}
