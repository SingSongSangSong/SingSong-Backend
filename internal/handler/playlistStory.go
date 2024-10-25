package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type KeepListResponse struct {
	KeepLists  []KeepLists `json:"keepLists"`
	LastCursor int64       `json:"lastCursor"`
}

type KeepLists struct {
	KeepListId int64            `json:"keepListId"`
	KeepName   string           `json:"keepName"`
	MemberId   int64            `json:"memberId"`
	Nickname   string           `json:"nickname"`
	LikeCount  int              `json:"likeCount"`
	UpdatedAt  time.Time        `json:"updatedAt"`
	IsLiked    bool             `json:"isLiked"`
	KeepSongs  []KeepStorySongs `json:"keepSongs"`
}

type KeepStorySongs struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	SongInfoId        int64  `json:"songId"`
	Album             string `json:"album"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
}

// GetKeepForStory godoc
// @Summary      최근 플레이리스트 List업 하기
// @Description  최근 플레이리스트를 여러 방법으로 가져오고 사용자의 좋아요 여부까지 알 수 있게 합니다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        filter query string false "정렬 기준. 최신순=recent, 오래된순(디폴트)=old, 좋아요가 많은순=like"
// @Param        size query string false "한번에 조회할 플레이리스트 개수. 디폴트값은 10 + @(노래개수)"
// @Param        cursor query string false "마지막에 조회했던 커서의 keep_list_id(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 정렬기준의 가장 처음 플리부터 줌"
// @Success      200 {object} pkg.BaseResponseStruct{data=KeepListResponse} "성공"
// @Router       /v1/keep/story [get]
// @Security BearerAuth
func GetKeepForStory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 기본은 최신순으로 정렬한다
		filter := c.DefaultQuery("filter", "recent")
		// filter가 recent, old, like가 아닌 경우 에러 반환
		if filter == "" || (filter != "recent" && filter != "old" && filter != "like") {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid filter in query", nil)
			return
		}

		// size와 cursor 파라미터를 받아온다
		sizeStr := c.DefaultQuery("size", "10")
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		// cursor가 없거나 0보다 작은 경우 에러 반환
		cursorStr := c.DefaultQuery("cursor", "9223372036854775807")
		if filter == "old" {
			cursorStr = c.DefaultQuery("cursor", "0")
		}
		cursorInt, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || cursorInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		// 차단된 회원 정보 조회
		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		// blockedMemberIds가 빈 슬라이스가 아닐 경우 쿼리에서 사용하기 위해 변환
		if len(blacklists) > 0 {
			for _, blacklist := range blacklists {
				blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
			}
		}

		// 기본 orderBy 및 cursorCondition 설정
		orderBy := "kl.updated_at DESC, kl.keep_list_id DESC"
		cursorCondition := "kl.keep_list_id < ?"
		cursorValue := cursorInt

		// filter 조건에 따른 orderBy와 cursorCondition 설정
		if filter == "old" {
			orderBy = "kl.updated_at ASC, kl.keep_list_id ASC"
			cursorCondition = "kl.keep_list_id > ?"
		} else if filter == "like" {
			orderBy = "kl.likes DESC, kl.updated_at DESC, kl.keep_list_id DESC"
			cursorCondition = "kl.keep_list_id < ?"
		}

		// 기본 SQL 쿼리 생성 각 KeepList별로 최근에 추가한 노래 4개만 보여주기
		query := `
			SELECT kl.keep_list_id, kl.keep_name, m.member_id, m.nickname, kl.likes, kl.updated_at, 
				(
					SELECT GROUP_CONCAT(ks2.song_info_id ORDER BY ks2.created_at DESC LIMIT 4)
					FROM keep_song AS ks2
					WHERE ks2.keep_list_id = kl.keep_list_id
				) AS songInfoIds,
				(
					SELECT COUNT(ks3.song_info_id)
					FROM keep_song AS ks3
					WHERE ks3.keep_list_id = kl.keep_list_id
				) AS songCount
			FROM keep_list AS kl
			LEFT JOIN member AS m ON kl.member_id = m.member_id
			WHERE kl.keep_name NOT LIKE '%Anonymous의 플레이리스트%'
			`

		// 블록된 회원 제외 조건 추가
		if len(blockedMemberIds) > 0 {
			query += "AND kl.member_id NOT IN (?) "
		}

		// cursor 조건 추가
		query += "AND " + cursorCondition + " "

		// group by 및 order by, limit 추가
		query += `HAVING songCount >= 3 ORDER BY ` + orderBy + " LIMIT ?"

		// SQL 쿼리 실행을 위한 파라미터 준비
		args := []interface{}{cursorValue, sizeInt}
		if len(blockedMemberIds) > 0 {
			args = append(blockedMemberIds, args...)
		}

		// Query를 통해 SQL 실행
		rows, err := db.Query(query, args...)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close()

		// 결과를 담을 구조체 슬라이스 생성
		keepListIds := make([]int64, 0, sizeInt)
		keepListsResponse := make([]KeepLists, 0, sizeInt)
		keepListSongMap := make(map[int64][]int64)
		uniqueSongIdsMap := make(map[int64]struct{}) // 중복을 방지하기 위한 map
		uniqueSongIds := make([]int64, 0)            // 고유한 ID를 담을 슬라이스
		// 조회 결과를 반복하면서 값을 스캔
		for rows.Next() {
			var keepLists KeepLists
			var songInfoIds null.String
			var songCount int
			err := rows.Scan(
				&keepLists.KeepListId,
				&keepLists.KeepName,
				&keepLists.MemberId,
				&keepLists.Nickname,
				&keepLists.LikeCount,
				&keepLists.UpdatedAt,
				&songInfoIds,
				&songCount,
			)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			keepListIds = append(keepListIds, keepLists.KeepListId)
			keepListsResponse = append(keepListsResponse, keepLists)

			songInfoList, err := ConvertToInt64List(songInfoIds.String)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			// songInfoList에서 uniqueSongIdsMap에 없는 값만 추가
			for _, songId := range songInfoList {
				if _, exists := uniqueSongIdsMap[songId]; !exists {
					uniqueSongIdsMap[songId] = struct{}{}         // 중복 방지를 위해 map에 추가
					uniqueSongIds = append(uniqueSongIds, songId) // 중복이 없으므로 uniqueSongIds에 추가
				}
			}
			// KeepListSongMap에 곡리스트 정보 추가
			keepListSongMap[keepLists.KeepListId] = songInfoList
		}

		// KeepListLikeMap을 생성 (좋아요 여부를 저장하는 맵)
		keepListLikeMap := make(map[int64]bool)
		keepMemberLikes, err := mysql.KeepListLikes(
			qm.Where("member_id = ?", blockerId),
			qm.WhereIn("keep_list_id IN ?", int64SliceToInterface(keepListIds)...),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		for _, like := range keepMemberLikes {
			keepListLikeMap[like.KeepListID] = true
		}

		songInfoList := make([]*mysql.SongInfo, 0)
		// songInfoId로 곡 정보를 조회
		if len(uniqueSongIds) > 0 {
			songInfoList, err = mysql.SongInfos(
				qm.WhereIn("song_info_id IN ?", int64SliceToInterface(uniqueSongIds)...),
			).All(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}

		// 곡 정보를 map에 저장
		songInfoMap := make(map[int64]KeepStorySongs)
		for _, songInfo := range songInfoList {
			songInfoMap[songInfo.SongInfoID] = KeepStorySongs{
				SongNumber:        songInfo.SongNumber,
				SongName:          songInfo.SongName,
				SingerName:        songInfo.ArtistName,
				SongInfoId:        songInfo.SongInfoID,
				Album:             songInfo.Album.String,
				IsMr:              songInfo.IsMR.Bool,
				IsLive:            songInfo.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(songInfo.MelonSongID),
				LyricsYoutubeLink: songInfo.LyricsVideoLink.String,
				TJYoutubeLink:     songInfo.TJYoutubeLink.String,
			}
		}

		// KeepListSongMap을 기반으로 곡 정보를 추가
		for i, keepLists := range keepListsResponse {
			for _, songInfoId := range keepListSongMap[keepLists.KeepListId] {
				keepListsResponse[i].KeepSongs = append(keepListsResponse[i].KeepSongs, songInfoMap[songInfoId])
			}
			keepListsResponse[i].IsLiked = keepListLikeMap[keepLists.KeepListId]
		}

		// 커서 설정
		lastCursor := int64(0)
		if len(keepListsResponse) > 0 {
			lastCursor = keepListsResponse[len(keepListsResponse)-1].KeepListId
		} else {
			lastCursor = cursorInt
		}

		// 결과 반환
		pkg.BaseResponse(c, http.StatusOK, "success", KeepListResponse{KeepLists: keepListsResponse, LastCursor: lastCursor})
	}
}

func ConvertToInt64List(songInfoIds string) ([]int64, error) {
	// 문자열이 빈 경우 바로 반환
	if songInfoIds == "" {
		return []int64{}, nil
	}

	// 문자열에서 대괄호 [] 제거
	songInfoIds = strings.ReplaceAll(songInfoIds, "[", "")
	songInfoIds = strings.ReplaceAll(songInfoIds, "]", "")

	// 문자열을 쉼표로 나누어 배열로 변환
	songInfoList := strings.Split(songInfoIds, ",")

	// 결과를 담을 int64 슬라이스 생성
	int64List := make([]int64, len(songInfoList))

	// 각 문자열을 int64로 변환
	for i, idStr := range songInfoList {
		id, err := strconv.ParseInt(idStr, 10, 64) // 10진수, 64비트로 변환
		if err != nil {
			return nil, err // 변환 중 오류가 발생하면 오류 반환
		}
		int64List[i] = id
	}

	return int64List, nil
}

// Convert int64 slice to []interface{} for passing into SQL queries
func int64SliceToInterface(slice []int64) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

// KeepListLike godoc
// @Summary      플레이리스트 좋아요/좋아요 취소
// @Description  플레이리스트에 좋아요를 누르거나 좋아요를 취소합니다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        keepListId path string true "keepListId"
// @Success      200 {object} pkg.BaseResponseStruct{data=nil} "성공"
// @Router       /v1/keep/{keepListId}/like [post]
// @Security BearerAuth
func KeepListLike(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// commentId 가져오기
		keepListIdParam := c.Param("keepListId")
		keepListId, err := strconv.ParseInt(keepListIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}

		// 좋아요 상태 변경 함수
		changeLikeStatus := func(keepList *mysql.KeepList, delta int) error {
			if keepList.Likes.Valid {
				keepList.Likes.Int += delta
			} else {
				keepList.Likes = null.IntFrom(delta)
			}
			_, err := keepList.Update(c, db, boil.Infer())
			return err
		}

		// 이미 좋아요를 눌렀는지 확인
		keepListLikes, err := mysql.KeepListLikes(
			qm.Where("member_id = ? AND keep_list_id = ? AND deleted_at IS NULL", memberId.(int64), keepListId),
		).One(c.Request.Context(), db)

		// 이미 좋아요를 누른 상태에서 좋아요 취소 요청
		if err == nil {
			keepListLikes.DeletedAt = null.TimeFrom(time.Now())
			if _, err := keepListLikes.Update(c.Request.Context(), db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// CommentTable에서 해당 CommentId의 LikeCount를 1 감소시킨다
			keepList, err := mysql.KeepLists(
				qm.Where("keep_list_id = ?", keepListId),
			).One(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			if err := changeLikeStatus(keepList, -1); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", nil)
			return
		}

		// keepList 좋아요 누르기
		like := mysql.KeepListLike{MemberID: memberId.(int64), KeepListID: keepListId}
		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// CommentTable에서 해당 CommentId의 LikeCount를 1 증가시킨다
		keepList, err := mysql.KeepLists(
			qm.Where("keep_list_id = ?", keepListId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if err := changeLikeStatus(keepList, 1); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
		return
	}
}

type GetSongsFromKeepInStoryResponse struct {
	Songs     []GetSongsFromKeepStories `json:"songs"`
	MemberId  int64                     `json:"memberId"`
	KeepName  string                    `json:"keepName"`
	LikeCount int                       `json:"likeCount"`
	IsLiked   bool                      `json:"isLiked"`
}

type GetSongsFromKeepStories struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	SongInfoId        int64  `json:"songId"`
	Album             string `json:"album"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	MelonLink         string `json:"melonLink"`
	LyricsYoutubeLink string `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string `json:"tjYoutubeLink"`
}

// GetSongsFromKeepInStory godoc
// @Summary      다른사람의 플레이리스트에 있는 노래를 가져온다
// @Description  다른사람의 플레이리스트에 있는 노래를 가져온다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        keepListId path string true "keepListId"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetSongsFromKeepInStoryResponse} "성공"
// @Router       /v1/keep/{keepListId} [get]
// @Security BearerAuth
func GetSongsFromKeepInStory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// keepListId 가져오기
		keepListIdParam := c.Param("keepListId")
		keepListId, err := strconv.ParseInt(keepListIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}

		// KeepList정보 가져오기
		m := mysql.KeepLists(qm.Where("keep_list_id = ?", keepListId))
		playlistInfo, errors := m.One(c.Request.Context(), db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+errors.Error(), nil)
			return
		}

		result := mysql.KeepSongs(qm.Where("keep_list_id = ? AND deleted_at IS NULL", playlistInfo.KeepListID))
		all, err2 := result.All(c.Request.Context(), db)
		if err2 != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err2.Error(), nil)
			return
		}

		// all을 songInfoIds로 바꾸기
		songInfoIds := make([]int64, len(all))
		for i, v := range all {
			songInfoIds[i] = v.SongInfoID
		}
		// songInfoIds로 SongInfos를 한번에 가져오기
		allSongInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", int64SliceToInterface(songInfoIds)...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// allSongInfos를 map으로 바꾸기
		songFromKeepResponse := make([]GetSongsFromKeepStories, len(allSongInfos))
		for i, v := range allSongInfos {
			songFromKeepResponse[i] = GetSongsFromKeepStories{
				SongName:          v.SongName,
				SingerName:        v.ArtistName,
				SongNumber:        v.SongNumber,
				SongInfoId:        v.SongInfoID,
				Album:             v.Album.String,
				IsMr:              v.IsMR.Bool,
				IsLive:            v.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(v.MelonSongID),
				LyricsYoutubeLink: v.LyricsVideoLink.String,
				TJYoutubeLink:     v.TJYoutubeLink.String,
			}
		}
		// 좋아요 여부 확인
		_, err = mysql.KeepListLikes(
			qm.Where("member_id = ? AND keep_list_id = ? AND deleted_at IS NULL", memberId, playlistInfo.KeepListID),
		).One(c.Request.Context(), db)
		if err == sql.ErrNoRows {
			// 결과 반환
			pkg.BaseResponse(c, http.StatusOK, "success", GetSongsFromKeepInStoryResponse{
				Songs:     songFromKeepResponse,
				MemberId:  playlistInfo.MemberID,
				KeepName:  playlistInfo.KeepName.String,
				LikeCount: playlistInfo.Likes.Int,
				IsLiked:   false,
			})
			return
		}
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 결과 반환
		pkg.BaseResponse(c, http.StatusOK, "success", GetSongsFromKeepInStoryResponse{
			Songs:     songFromKeepResponse,
			MemberId:  playlistInfo.MemberID,
			KeepName:  playlistInfo.KeepName.String,
			LikeCount: playlistInfo.Likes.Int,
			IsLiked:   true,
		})
		return
	}
}

// SubscribeKeep godoc
// @Summary      플레이리스트 구독하기/구독 취소
// @Description  플레이리스트를 구독하거나 구독을 취소합니다. 요청을 한번 보내면 구독을 하고 한번 더 보내면 구독 여부 확인후 취소합니다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        keepListId path string true "keepListId"
// @Success      200 {object} pkg.BaseResponseStruct{data=nil} "성공"
// @Router       /v1/keep/{keepListId}/subscribe [post]
// @Security BearerAuth
func SubscribeKeep(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// memberId Int64로 변환
		memberIdInt, ok := memberId.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - invalid memberId type", nil)
			return
		}

		// keepListId 가져오기
		keepListIdParam := c.Param("keepListId")
		keepListId, err := strconv.ParseInt(keepListIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}

		// 이미 구독했는지 확인
		keepListSubscribes, err := mysql.KeepListSubscribes(
			qm.Where("member_id = ? AND keep_list_id = ? AND deleted_at IS NULL", memberIdInt, keepListId),
		).One(c.Request.Context(), db)
		if err == nil {
			// 이미 구독한 경우 구독을 취소해야한다
			keepListSubscribes.DeletedAt = null.TimeFrom(time.Now())
			if _, err := keepListSubscribes.Update(c.Request.Context(), db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			pkg.BaseResponse(c, http.StatusOK, "success", nil)
			return
		}

		// 구독하지 않은 경우 구독한다
		subscribe := mysql.KeepListSubscribe{MemberID: memberIdInt, KeepListID: keepListId}
		if err := subscribe.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

type GetSubscribedKeepsResponse struct {
	KeepListId int64     `json:"keepListId"`
	KeepName   string    `json:"keepName"`
	MemberId   int64     `json:"memberId"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// GetSubscribedKeeps godoc
// @Summary      내가 구독한 다른사람들의 플레이리스트 가져오기
// @Description  내가 구독한 다른사람들의 플레이리스트를 가져온다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]GetSubscribedKeepsResponse} "성공"
// @Router       /v1/keep/subscribe [get]
// @Security BearerAuth
func GetSubscribedKeeps(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// 구독한 플레이리스트 정보 가져오기
		subscribedKeeps, err := mysql.KeepListSubscribes(
			qm.Where("member_id = ? AND deleted_at IS NULL", memberId),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 구독한 플레이리스트 정보를 담을 interface 생성
		subscribedKeepListIds := make([]int64, len(subscribedKeeps))
		for i, v := range subscribedKeeps {
			subscribedKeepListIds[i] = v.KeepListID
		}

		// KeepList정보 가져오기
		m := mysql.KeepLists(
			qm.WhereIn("keep_list_id In ?", int64SliceToInterface(subscribedKeepListIds)...),
		)
		playlistInfo, errors := m.All(c.Request.Context(), db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+errors.Error(), nil)
			return
		}

		// 리턴 객체 반환
		SubscribedKeepsResponseList := make([]GetSubscribedKeepsResponse, len(playlistInfo))
		for i, v := range playlistInfo {
			SubscribedKeepsResponseList[i] = GetSubscribedKeepsResponse{
				KeepListId: v.KeepListID,
				KeepName:   v.KeepName.String,
				MemberId:   v.MemberID,
				UpdatedAt:  v.UpdatedAt.Time,
			}
		}

		// 결과 반환
		pkg.BaseResponse(c, http.StatusOK, "success", SubscribedKeepsResponseList)

	}
}
