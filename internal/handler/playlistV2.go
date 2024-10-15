package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
)

type GetPlayListV2Response struct {
	PlayListResponse []PlaylistAddResponse `json:"songs"`
	LastCursor       int64                 `json:"lastCursor"`
}

type KeepSongsWithSongName struct {
	SongNumber  int         `json:"songNumber"`
	SongName    string      `json:"songName"`
	ArtistName  string      `json:"artistName"`
	SongInfoId  int64       `json:"songInfoId"`
	Album       null.String `json:"album"`
	IsMr        null.Bool   `json:"isMr"`
	IsLive      null.Bool   `json:"isLive"`
	MelonSongId null.String `json:"melonSongId"`
}

// GetSongsFromPlaylistV2 godoc
// @Summary      플레이리스트에 노래를 여러가지 필터로 가져온다 (커서기반 페이징)
// @Description  플레이리스트에 있는 노래들을 가나다순/최신추가순/오래된순(alphabet/recent/old) 으로 가져온다. 기본값은 recent이다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        filter query string false "필터"
// @Param        cursor query int false "마지막에 조회했던 커서의 songId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신 글부터 조회"
// @Param        size query int false "한번에 조회할 노래의 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetPlayListV2Response} "성공"
// @Router       /v2/keep [get]
// @Security BearerAuth
func GetSongsFromPlaylistV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, err1 := c.Get("memberId")
		if err1 != true {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// 검색어를 쿼리 파라미터에서 가져오기
		filter := c.DefaultQuery("filter", "recent")
		if filter == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find filter in query", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSearchSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		// Playlist정보 가져오기
		m := mysql.KeepLists(qm.Where("member_id = ?", memberId))
		playlistInfo, errors := m.One(c.Request.Context(), db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+errors.Error(), nil)
			return
		}

		// 기본 커서 값 및 정렬 기준 설정
		var orderClause, cursorCondition string
		var cursorInt int64
		cursorStr := c.DefaultQuery("cursor", "0") // 기본값으로 0
		cursorInt, err = strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || cursorInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		// 필터에 따른 정렬 및 커서 처리
		switch filter {
		case "alphabet":
			orderClause = "ORDER BY song_info.song_name ASC"
			cursorCondition = "AND song_info.song_info_id > ?"
		case "recent":
			cursorStr = c.DefaultQuery("cursor", "9223372036854775807") // int64 최대값
			cursorInt, err = strconv.ParseInt(cursorStr, 10, 64)
			if err != nil || cursorInt < 0 {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
				return
			}
			orderClause = "ORDER BY keep_song.keep_song_id DESC"
			cursorCondition = "AND keep_song.keep_song_id < ?"
		case "old":
			orderClause = "ORDER BY keep_song.keep_song_id ASC"
			cursorCondition = "AND keep_song.keep_song_id > ?"
		default:
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid filter parameter", nil)
			return
		}

		// 공통 쿼리 생성
		query := fmt.Sprintf(`
			SELECT song_info.song_number, song_info.song_name, song_info.artist_name, song_info.song_info_id, song_info.album, song_info.is_mr, song_info.is_live, song_info.melon_song_id, keep_song.keep_song_id
			FROM keep_song
			LEFT JOIN song_info ON keep_song.song_info_id = song_info.song_info_id
			WHERE keep_song.keep_list_id = ? AND keep_song.deleted_at IS NULL %s %s
			LIMIT ?
		`, cursorCondition, orderClause)

		// Query 실행
		rows, err := db.Query(query, playlistInfo.KeepListID, cursorInt, sizeInt)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close()

		// 결과를 담을 구조체 슬라이스 생성
		keepSongs := make([]PlaylistAddResponse, 0, sizeInt)

		// 조회 결과를 반복하면서 값을 스캔
		for rows.Next() {
			var keepSong KeepSongsWithSongName
			var keepSongId int64
			err := rows.Scan(
				&keepSong.SongNumber,
				&keepSong.SongName,
				&keepSong.ArtistName,
				&keepSong.SongInfoId,
				&keepSong.Album,
				&keepSong.IsMr,
				&keepSong.IsLive,
				&keepSong.MelonSongId,
				&keepSongId,
			)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			playlistAddResponse := PlaylistAddResponse{
				SongNumber: keepSong.SongNumber,
				SongName:   keepSong.SongName,
				SingerName: keepSong.ArtistName,
				SongInfoId: keepSong.SongInfoId,
				Album:      keepSong.Album.String,
				IsMr:       keepSong.IsMr.Bool,
				IsLive:     keepSong.IsLive.Bool,
				MelonLink:  CreateMelonLinkByMelonSongId(null.StringFrom(keepSong.MelonSongId.String)),
				KeepSongId: keepSongId,
			}
			keepSongs = append(keepSongs, playlistAddResponse)
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(keepSongs) > 0 && filter == "alphabet" {
			lastCursor = keepSongs[len(keepSongs)-1].SongInfoId
		} else if len(keepSongs) > 0 && (filter == "recent" || filter == "old") {
			lastCursor = keepSongs[len(keepSongs)-1].KeepSongId
		}
		if len(keepSongs) == 0 && filter == "alphabet" || filter == "old" {
			lastCursor = cursorInt
		}
		getPlayListV2Response := GetPlayListV2Response{
			PlayListResponse: keepSongs,
			LastCursor:       lastCursor,
		}

		// 성공적으로 처리된 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "success", getPlayListV2Response)
	}
}
