package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	pb "SingSong-Server/proto/functionCallingRecommend"
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
)

type FunctionCallingDetailResponse struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	SongInfoId        int64  `json:"songId"`
	Album             string `json:"album"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	IsKeep            bool   `json:"isKeep"`
	KeepCount         int    `json:"keepCount"`
	CommentCount      int    `json:"commentCount"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
	LyricsVideoID     string `json:"lyricsVideoId"`
	TJVideoID         string `json:"tjVideoId"`
}

type FunctionCallingResponse struct {
	Songs []FunctionCallingDetailResponse `json:"songs"`
}

// FunctionCallingRecommedation godoc
// @Summary      LLM으로 검색하기
// @Description  LLM의 사용자 입력을 토대로 추천된 노래를 반환합니다. 10개의 노래를 반환합니다
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        input   body      LlmRequest  true  "인풋"
// @Success      200 {object} pkg.BaseResponseStruct{data=FunctionCallingResponse} "성공"
// @Router       /v1/recommend/recommendation/functionCalling [post]
// @Security BearerAuth
func FunctionCallingRecommedation(db *sql.DB) gin.HandlerFunc {
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
			log.Printf("Did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewFunctionCallingRecommendClient(conn)

		// gRPC 요청 생성
		rpcRequest := &pb.FunctionCallingRequest{
			MemberId: memberIdInt,
			Command:  llmRequest.UserInput,
		}

		go func() {
			llmSearch := mysql.LLMSearchLog{MemberID: memberIdInt, SearchText: llmRequest.UserInput}
			err = llmSearch.Insert(c.Request.Context(), db, boil.Infer())
			if err != nil {
				log.Printf("Error inserting LLM Search Log: %v", err)
			}
		}()

		// gRPC 요청 보내기
		response, err := client.GetFunctionCallingRecommendation(context.Background(), rpcRequest)
		if err != nil {
			log.Printf("Error calling gRPC: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Populate the userProfileResponse with gRPC response data
		functionCallingResponse := FunctionCallingResponse{
			Songs: []FunctionCallingDetailResponse{},
		}

		// SongInfoId 리스트를 담을 빈 리스트 생성
		var songInfoIds []int64

		// gRPC response에서 SongInfoId만 추출
		for _, item := range response.SongInfoId {
			songInfoIds = append(songInfoIds, item)
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

		// Song정보를 저장하는 맵 생성
		songInfoMap := make(map[int64]*mysql.SongInfo)
		for _, songInfo := range songInfos {
			songInfoMap[songInfo.SongInfoID] = songInfo
		}

		keepList, err := mysql.KeepLists(
			qm.Where("member_id = ?", memberIdInt),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongs, err := mysql.KeepSongs(
			qm.Where("keep_list_id = ?", keepList.KeepListID),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepSongMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepSongMap[keepSong.SongInfoID] = true
		}

		// songInfoInterface 배열을 "?, ?, ?" 형식으로 변환
		inClause := strings.TrimSuffix(strings.Repeat("?,", len(songInfoInterface)), ",")

		// 쿼리 작성 (KeepCount)
		keepCountQuery := fmt.Sprintf(`
			SELECT keep_song.song_info_id, COUNT(keep_song.keep_song_id) AS song_count
			FROM keep_song
			WHERE keep_song.song_info_id IN (%s)
			AND keep_song.deleted_at IS NULL
			GROUP BY keep_song.song_info_id
		`, inClause)

		// 쿼리 실행 (KeepCount)
		rows, err := db.QueryContext(c.Request.Context(), keepCountQuery, songInfoInterface...)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close() // 쿼리 종료 후 닫음

		// 결과를 저장할 맵 생성 (KeepCount)
		keepCountMap := make(map[int64]int)
		for rows.Next() {
			var songInfoId int64
			var songCount int
			if err := rows.Scan(&songInfoId, &songCount); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			keepCountMap[songInfoId] = songCount
		}

		// commentCount Query 작성
		commentCountQuery := fmt.Sprintf(`
			SELECT comment.song_info_id, COUNT(comment_id) AS comment_count
			FROM comment
			WHERE comment.song_info_id IN (%s)
			AND comment.deleted_at IS NULL
			GROUP BY comment.song_info_id
		`, inClause)

		// 쿼리 실행 (CommentCount)
		commentRows, err := db.QueryContext(c.Request.Context(), commentCountQuery, songInfoInterface...)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer commentRows.Close() // 쿼리 종료 후 닫음

		// 결과를 저장할 맵 생성 (CommentCount)
		commentCountMap := make(map[int64]int)
		for commentRows.Next() {
			var songInfoId int64
			var commentCount int
			if err := commentRows.Scan(&songInfoId, &commentCount); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			commentCountMap[songInfoId] = commentCount
		}

		// Loop through the gRPC response to populate songResponse
		for _, item := range response.SongInfoId {
			functionCallingResponse.Songs = append(functionCallingResponse.Songs, FunctionCallingDetailResponse{
				SongNumber:        songInfoMap[item].SongNumber,
				SongName:          songInfoMap[item].SongName,
				SingerName:        songInfoMap[item].ArtistName,
				SongInfoId:        songInfoMap[item].SongInfoID,
				Album:             songInfoMap[item].Album.String,
				IsMr:              songInfoMap[item].IsMR.Bool,
				IsLive:            songInfoMap[item].IsLive.Bool,
				IsKeep:            keepSongMap[songInfoMap[item].SongInfoID],
				KeepCount:         keepCountMap[songInfoMap[item].SongInfoID],
				CommentCount:      commentCountMap[songInfoMap[item].SongInfoID],
				MelonLink:         CreateMelonLinkByMelonSongId(songInfoMap[item].MelonSongID),
				LyricsYoutubeLink: songInfoMap[item].LyricsVideoLink.String,
				TJYoutubeLink:     songInfoMap[item].TJYoutubeLink.String,
				LyricsVideoID:     ExtractVideoID(songInfoMap[item].LyricsVideoLink.String),
				TJVideoID:         ExtractVideoID(songInfoMap[item].TJYoutubeLink.String),
			})
		}

		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", functionCallingResponse)
	}
}
