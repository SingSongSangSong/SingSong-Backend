package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"time"
)

type newSongInfo struct {
	SongNumber        int    `json:"songNumber"`
	SongName          string `json:"songName"`
	SingerName        string `json:"singerName"`
	Album             string `json:"album"`
	IsKeep            bool   `json:"isKeep"`
	SongInfoId        int64  `json:"songId"`
	IsMr              bool   `json:"isMr"`
	IsLive            bool   `json:"isLive"`
	KeepCount         int    `json:"keepCount"`
	CommentCount      int    `json:"commentCount"`
	MelonLink         string `json:"melonLink"`
	IsRecentlyUpdated bool   `json:"isRecentlyUpdated"`
}

type newSongInfoResponse struct {
	Songs      []newSongInfo `json:"songs"`
	LastCursor int64         `json:"lastCursor"`
}

// ListNewSongs godoc
// @Summary      최근 한달간의 신곡을 조회 (최신순 조회)
// @Description  최근 한달간의 신곡을 최신순으로 조회합니다. 최근 일주일동안 추가된 신곡은 isRecentlyUpdated = true 입니다.
// @Tags         New Song
// @Accept       json
// @Produce      json
// @Param        cursor query int false "마지막에 조회했던 커서의 songId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신곡부터 조회"
// @Param        size query int false "한번에 조회할 노래 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=newSongInfoResponse} "성공"
// @Failure      400 "query param 값이 들어왔는데, 숫자가 아니라면 400 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/songs/new [get]
// @Security BearerAuth
func ListNewSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt <= 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		// 페이징 처리된 신곡 가져오기
		songInfos, err := mysql.SongInfos(
			qm.Where("song_info_id < ?", cursorInt),
			qm.And("created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)"), //최근 31일동안 추가된 곡
			qm.OrderBy("song_info_id DESC"),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Keep 여부 가져오기
		keepLists, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepListInterface := make([]interface{}, len(keepLists))
		for i, v := range keepLists {
			keepListInterface[i] = v.KeepListID
		}
		// []int64를 []interface{}로 변환
		songInfoInterface := make([]interface{}, len(songInfos))
		for i, v := range songInfos {
			songInfoInterface[i] = v.SongInfoID
		}
		keepSongs, err := mysql.KeepSongs(
			qm.WhereIn("keep_list_id = ?", keepListInterface...),
			qm.AndIn("song_info_id IN ?", songInfoInterface...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepMap := make(map[int64]bool)
		for _, keepSong := range keepSongs {
			keepMap[keepSong.SongInfoID] = true
		}

		newSongs := make([]newSongInfo, 0, sizeInt)
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)

		for _, song := range songInfos {
			isKeep := keepMap[song.SongInfoID]

			newSongs = append(newSongs, newSongInfo{
				SongNumber:        song.SongNumber,
				SongName:          song.SongName,
				SingerName:        song.ArtistName,
				Album:             song.Album.String,
				IsKeep:            isKeep,
				SongInfoId:        song.SongInfoID,
				IsMr:              song.IsMR.Bool,
				IsLive:            song.IsLive.Bool,
				KeepCount:         0, //todo
				CommentCount:      0, //todo
				MelonLink:         CreateMelonLinkByMelonSongId(song.MelonSongID),
				IsRecentlyUpdated: song.CreatedAt.Time.After(sevenDaysAgo), //song.CreatedAt이 7일 이내인지
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(newSongs) > 0 {
			lastCursor = newSongs[len(newSongs)-1].SongInfoId
		}

		response := newSongInfoResponse{
			Songs:      newSongs,
			LastCursor: lastCursor,
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
