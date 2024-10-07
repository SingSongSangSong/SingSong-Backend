package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	pb "SingSong-Server/proto/userProfileRecommend"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

// songHomeResponse와 songResponse가 동일한 것으로 가정하고 사용
type songResponse struct {
	SongNumber   int    `json:"songNumber"`
	SongName     string `json:"songName"`
	SingerName   string `json:"singerName"`
	SongInfoId   int64  `json:"songId"`
	Album        string `json:"album"`
	IsMr         bool   `json:"isMr"`
	IsLive       bool   `json:"isLive"`
	IsKeep       bool   `json:"isKeep"`
	KeepCount    int    `json:"keepCount"`
	CommentCount int    `json:"commentCount"`
	MelonLink    string `json:"melonLink"`
}

type userProfileResponse struct {
	Songs []songResponse `json:"songs"`
}

var (
	GrpcAddr = conf.GrpcConfigInstance.Addr
)

// GetRecommendation godoc
// @Summary      AI가 골랐송
// @Description  사용자의 프로필을 기반으로 추천된 노래를 반환합니다. 페이지당 20개의 노래를 반환합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=userProfileResponse} "성공"
// @Router       /v1/recommend/recommendation/ai [get]
// @Router       /v1/recommend/recommendation/{pageId} [get]
// @Security BearerAuth
func GetRecommendation(db *sql.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		conn, err := grpc.Dial(GrpcAddr+":50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()

		client := pb.NewUserProfileClient(conn)

		historySongsForMilvus := getRefreshHistoryForMilvus(c, redisClient, memberIdInt)
		vectorQuerySize := pageSize + len(historySongsForMilvus)

		// gRPC 요청 생성
		rpcRequest := &pb.ProfileRequest{
			MemberId: memberIdInt,
			Page:     int32(vectorQuerySize),
			Gender:   gender.(string),
		}

		// gRPC 요청 보내기
		response, err := client.CreateUserProfile(context.Background(), rpcRequest)
		if err != nil {
			log.Printf("Error calling gRPC: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		querySongs := extractSongInfoForMilvus(vectorQuerySize, response)

		// 이전에 조회한 노래 빼고 상위 pageSize개 선택
		refreshedSongs := getTopSongsWithoutHistory(historySongsForMilvus, querySongs)

		// 무한 새로고침 - (페이지의 끝일 때/노래 개수가 애초에 PageSize수보다 작을때) 부족한 노래 수만큼 refreshedSongs를 마저 채운다
		if len(refreshedSongs) < pageSize {
			refreshedSongs = fillSongsAgain(refreshedSongs, querySongs)
			// 기록 비우기
			historySongsForMilvus = []int64{}
		}

		// Populate the userProfileResponse with gRPC response data
		userProfileRes := userProfileResponse{
			Songs: []songResponse{},
		}

		// SongInfoId 리스트를 담을 빈 리스트 생성
		var songInfoIds []int64

		// gRPC response에서 SongInfoId만 추출
		for _, item := range refreshedSongs {
			songInfoIds = append(songInfoIds, item.SongInfoId)
		}

		// []int64를 []interface{}로 변환
		songInfoInterface := make([]interface{}, len(songInfoIds))
		for i, v := range songInfoIds {
			songInfoInterface[i] = v
		}

		// Keep 여부 가져오기
		keepLists, err := mysql.KeepLists(qm.Where("member_id = ?", memberIdInt)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepListInterface := make([]interface{}, len(keepLists))
		for i, v := range keepLists {
			keepListInterface[i] = v
		}
		keepSongs, err := mysql.KeepSongs(
			qm.WhereIn("keep_list_id = ?", keepListInterface...),
			qm.AndIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
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

		// MelonSongId 가져오기
		songInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
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

		// MelonSongId를 저장하는 맵 생성
		songInfoMap := make(map[int64]*mysql.SongInfo)
		for _, songInfo := range songInfos {
			songInfoMap[songInfo.SongInfoID] = songInfo
		}

		// gRPC response에서 가져온 SongInfoId를 기반으로 songInfoMap, keepSongsMap, commentsCountsMap, keepCountsMap을 활용
		for _, item := range refreshedSongs {
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
				SongNumber:   songInfoMap[item.SongInfoId].SongNumber,
				SongName:     songInfoMap[item.SongInfoId].SongName,
				SingerName:   songInfoMap[item.SongInfoId].ArtistName,
				SongInfoId:   item.SongInfoId,
				Album:        songInfoMap[item.SongInfoId].Album.String,
				IsMr:         songInfoMap[item.SongInfoId].IsMR.Bool,
				IsLive:       songInfoMap[item.SongInfoId].IsLive.Bool,
				IsKeep:       isKeep,       // Keep 여부 추가
				CommentCount: commentCount, // 댓글 수 추가
				KeepCount:    keepCount,    // Keep 수 추가
				MelonLink:    CreateMelonLinkByMelonSongId(songInfoMap[item.SongInfoId].MelonSongID),
			})
		}

		// history 갱신
		for _, song := range refreshedSongs {
			historySongsForMilvus = append(historySongsForMilvus, song.SongInfoId)
		}
		setRefreshHistoryForMilvus(c, redisClient, memberIdInt, historySongsForMilvus)

		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", userProfileRes)
	}
}

func getRefreshHistoryForMilvus(c *gin.Context, redisClient *redis.Client, memberId int64) []int64 {
	key := generateRefreshKeyForMilvus(memberId)

	val, err := redisClient.Get(c, key).Result()
	if err == redis.Nil {
		return []int64{}
	} else if err != nil {
		log.Printf("Failed to get history from Redis: %v", err)
		return []int64{}
	}

	var history []int64
	err = json.Unmarshal([]byte(val), &history)
	if err != nil {
		log.Printf("Failed to unmarshal history: %v", err)
		return []int64{}
	}

	return history
}

func generateRefreshKeyForMilvus(memberId int64) string {
	return "aiRecommendation:" + strconv.FormatInt(memberId, 10) + ":" + "milvus"
}

func extractSongInfoForMilvus(vectorQuerySize int, values *pb.ProfileResponse) []refreshResponse {
	querySongs := make([]refreshResponse, 0, vectorQuerySize)
	for _, match := range values.GetSimilarItems() {
		querySongs = append(querySongs, refreshResponse{
			SongInfoId: match.SongInfoId,
		})
	}
	return querySongs
}

func setRefreshHistoryForMilvus(c *gin.Context, redisClient *redis.Client, memberId int64, history []int64) {
	key := generateRefreshKeyForMilvus(memberId)

	historyJSON, err := json.Marshal(history)
	if err != nil {
		log.Printf("Failed to marshal history: %v", err)
		return
	}

	redisClient.Set(c, key, historyJSON, 30*time.Minute)
}
