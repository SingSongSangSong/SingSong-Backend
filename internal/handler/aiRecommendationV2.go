package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"fmt"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus-proto/go-api/v2/schemapb"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/null/v8"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SearchParam Custom search parameter struct that implements SearchParam interface
type SearchParam struct {
	params map[string]interface{}
}

// Params Ensure SearchParam satisfies the SearchParam interface
func (sp *SearchParam) Params() map[string]interface{} {
	return sp.params
}

func (sp *SearchParam) AddRadius(radius float64) {
	sp.params["radius"] = radius
}

func (sp *SearchParam) AddRangeFilter(rangeFilter float64) {
	sp.params["range_filter"] = rangeFilter
}

// NewSearchParam creates a new SearchParam with default parameters
func NewSearchParam() *SearchParam {
	return &SearchParam{
		params: map[string]interface{}{
			"nprobe": 10, // Example nprobe value for search granularity
		},
	}
}

// GetRecommendationV2 godoc
// @Summary      AI가 골랐송 Without GRPC
// @Description  사용자의 프로필을 기반으로 추천된 노래를 반환합니다. 페이지당 20개의 노래를 반환합니다.
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        pageId path int true "Page ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=UserProfileResponse} "성공"
// @Router       /v2/recommend/recommendation/ai [get]
// @Router       /v2/recommend/recommendation/{pageId} [get]
// @Security BearerAuth
func GetRecommendationV2(db *sql.DB, redisClient *redis.Client, milvus *client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 및 gender 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		gender, exists := c.Get("gender")
		if !exists {
			gender = "UNKNOWN"
		}

		// memberId를 int64로 변환
		memberIdInt, ok := memberId.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - invalid memberId type", nil)
			return
		}

		historySongsForMilvus := getRefreshHistoryForMilvus(c, redisClient, memberIdInt)
		vectorQuerySize := pageSize + len(historySongsForMilvus)

		// 유저 프로필 벡터 가져오기
		userVector, err := getUserProfile(milvus, memberIdInt)
		if err != nil || userVector == nil || len(userVector.GetFloatVector().Data) == 0 {
			log.Printf("No profile found for memberId %d fetching Gender Profile", memberIdInt)
			userVector, err = fetchGenderProfile(milvus, gender.(string))
			if err != nil || userVector == nil || len(userVector.GetFloatVector().Data) == 0 {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - no user profile vector", nil)
				return
			}
		}

		// 추천 노래 가져오기
		recommendation, err := recommendSimilarSongs(milvus, userVector, vectorQuerySize)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		querySongs := extractSongInfoForMilvusV2(vectorQuerySize, recommendation)
		// 이전에 조회한 노래 빼고 상위 pageSize개 선택
		refreshedSongs := getTopSongsWithoutHistory(historySongsForMilvus, querySongs)

		// 무한 새로고침 - (페이지의 끝일 때/노래 개수가 애초에 PageSize수보다 작을때) 부족한 노래 수만큼 refreshedSongs를 마저 채운다
		if len(refreshedSongs) < pageSize {
			refreshedSongs = fillSongsAgain(refreshedSongs, querySongs)
			// 기록 비우기
			historySongsForMilvus = []int64{}
		}

		// SongInfoId 리스트를 담을 빈 리스트 생성
		var songInfoIds []int64

		// 조회결과에서 SongInfoId만 추출
		for _, item := range refreshedSongs {
			songInfoIds = append(songInfoIds, item.SongInfoId)
		}

		// []int64를 []interface{}로 변환
		songInfoInterface := make([]interface{}, len(songInfoIds))
		for i, v := range songInfoIds {
			songInfoInterface[i] = v
		}

		// SQL 실행 전에 songInfoIds가 비어 있는지 확인
		if len(songInfoIds) == 0 {
			pkg.BaseResponse(c, http.StatusInternalServerError, "err - no song Ids", nil)
			return
		}

		// IN 절에 사용할 플레이스홀더 생성 (예: ?, ?, ?...)
		placeholders := make([]string, len(songInfoInterface))
		for i := range songInfoInterface {
			placeholders[i] = "?"
		}
		placeholderStr := strings.Join(placeholders, ", ")

		// SQL 쿼리: 필요한 데이터들을 JOIN을 통해 한 번에 가져오기
		query := fmt.Sprintf(`
			SELECT 
				si.song_info_id, si.song_number, si.song_name, si.artist_name, 
				si.album, si.is_mr, si.is_live, si.melon_song_id, si.lyrics_video_link, si.tj_youtube_link,
				COUNT(DISTINCT c.comment_id) AS comment_count,
				COUNT(DISTINCT ks.keep_song_id) AS keep_count,
				EXISTS (
					SELECT 1 FROM keep_song WHERE song_info_id = si.song_info_id AND keep_list_id IN (
						SELECT keep_list_id FROM keep_list WHERE member_id = ?
					)
				) AS is_keep
			FROM song_info si
			LEFT JOIN comment c ON si.song_info_id = c.song_info_id
			LEFT JOIN keep_song ks ON si.song_info_id = ks.song_info_id
			WHERE si.song_info_id IN (%s)
			GROUP BY si.song_info_id
		`, placeholderStr)

		// 쿼리 실행에 사용할 매개변수 준비 (memberIdInt + songInfoIds)
		args := append([]interface{}{memberIdInt}, songInfoInterface...)

		// SQL 쿼리 실행 및 결과 매핑
		rows, err := db.Query(query, args...)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close()

		// 결과를 SongResponse로 매핑
		var songs []SongResponse
		for rows.Next() {
			var song SongResponse
			var melonLinkId null.String
			var album null.String
			var lyricsLink null.String
			var tjLink null.String

			err := rows.Scan(
				&song.SongInfoId, &song.SongNumber, &song.SongName, &song.SingerName,
				&album, &song.IsMr, &song.IsLive, &melonLinkId, &lyricsLink, &tjLink,
				&song.CommentCount, &song.KeepCount, &song.IsKeep,
			)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			if album.Valid {
				song.Album = album.String
			}

			song.MelonLink = CreateMelonLinkByMelonSongId(melonLinkId)
			songs = append(songs, song)
		}

		// history 갱신
		for _, song := range refreshedSongs {
			historySongsForMilvus = append(historySongsForMilvus, song.SongInfoId)
		}
		setRefreshHistoryForMilvus(c, redisClient, memberIdInt, historySongsForMilvus)

		// 최종 응답 생성
		userProfileRes := UserProfileResponse{Songs: songs}
		// 결과를 JSON 형식으로 반환
		pkg.BaseResponse(c, http.StatusOK, "success", userProfileRes)
	}
}

// getUserProfile: Milvus 유저 프로필 벡터 조회
func getUserProfile(milvus *client.Client, memberID int64) (*schemapb.VectorField, error) {
	expr := "member_id == " + strconv.FormatInt(memberID, 10)
	res, err := (*milvus).Query(
		context.Background(),
		"user_profile",
		[]string{},
		expr,
		[]string{"profile_vector"},
	)
	if err != nil || len(res) == 0 {
		return nil, errors.New("no profile found")
	}
	//dimension := res.GetColumn("profile_vector").FieldData().GetVectors().Dim
	profileColumn := res.GetColumn("profile_vector").FieldData().GetVectors().GetFloatVector().Data
	if len(profileColumn) == 0 {
		return nil, errors.New("no profile found")
	}
	return res.GetColumn("profile_vector").FieldData().GetVectors(), nil
}

// fetchGenderProfile: 성별 기반 기본 프로필 조회
func fetchGenderProfile(milvus *client.Client, gender string) (*schemapb.VectorField, error) {
	var memberID int64
	switch gender {
	case "MALE":
		memberID = 0
	case "FEMALE":
		memberID = -1
	default:
		memberID = -2
	}
	return getUserProfile(milvus, memberID)
}

// recommendSimilarSongs: 유사도 기반 추천
func recommendSimilarSongs(milvus *client.Client, userVector *schemapb.VectorField, topK int) ([]int64, error) {
	userVectors := entity.FloatVector(userVector.GetFloatVector().Data)
	// Initialize search parameters
	searchParams := NewSearchParam()

	// Execute the search
	res, err := (*milvus).Search(
		context.Background(),
		conf.VectorDBConfigInstance.COLLECTION_NAME, // Collection name
		[]string{},
		"", // Partitions (empty for all partitions)
		[]string{"song_info_id", "song_name", "artist_name"}, // Output fields
		[]entity.Vector{userVectors},                         // Wrap the vector in a slice
		"vector",
		entity.COSINE, // Vector field name
		topK,          // Top-K results
		searchParams,  // Search parameters
	)
	if err != nil {
		log.Printf("Search failed: %v", err)
		return nil, err
	}

	return res[0].IDs.FieldData().GetScalars().GetLongData().Data, nil
}

func extractSongInfoForMilvusV2(vectorQuerySize int, values []int64) []refreshResponse {
	querySongs := make([]refreshResponse, 0, vectorQuerySize)
	for _, match := range values {
		querySongs = append(querySongs, refreshResponse{
			SongInfoId: match,
		})
	}
	return querySongs
}
