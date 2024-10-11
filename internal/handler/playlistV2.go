package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
)

type GetPlayListV2Response struct {
	PlayListResponse []PlaylistAddResponse `json:"playListResponse"`
	LastCursor       int64                 `json:"lastCursor"`
}

// GetSongsFromPlaylistV2 godoc
// @Summary      플레이리스트에 노래를 여러가지 필터로 가져온다 (커서기반 페이징)
// @Description  플레이리스트에 있는 노래들을 가나다순/최신추가순/오래된순(alphabet/recent/old) 으로 가져온다. 기본값은 recent이다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        filter query string true "필터"
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

		var all mysql.KeepSongSlice
		var err2 error

		if filter == "alphabet" {
			cursorStr := c.DefaultQuery("cursor", "0") //int64 최대값
			cursorInt, err := strconv.Atoi(cursorStr)
			if err != nil || cursorInt <= 0 {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
				return
			}

			result := mysql.KeepSongs(
				qm.Where("keep_list_id = ? AND deleted_at IS NULL AND song_info_id > ?", playlistInfo.KeepListID, cursorInt),
				qm.OrderBy("song_name ASC"),
				qm.Limit(sizeInt),
			)
			all, err2 = result.All(c.Request.Context(), db)
			if err2 != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err2.Error(), nil)
				return
			}
		} else if filter == "recent" {
			cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
			cursorInt, err := strconv.Atoi(cursorStr)
			if err != nil || cursorInt <= 0 {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
				return
			}

			result := mysql.KeepSongs(
				qm.Where("keep_list_id = ? AND deleted_at IS NULL AND song_info_id < ?", playlistInfo.KeepListID, cursorInt),
				qm.OrderBy("created_at DESC"),
				qm.Limit(sizeInt),
			)
			all, err2 = result.All(c.Request.Context(), db)
			if err2 != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err2.Error(), nil)
				return
			}
		} else if filter == "old" {
			cursorStr := c.DefaultQuery("cursor", "0") //int64 최대값
			cursorInt, err := strconv.Atoi(cursorStr)
			if err != nil || cursorInt <= 0 {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
				return
			}

			result := mysql.KeepSongs(
				qm.Where("keep_list_id = ? AND deleted_at IS NULL AND song_info_id > ?", playlistInfo.KeepListID, cursorInt),
				qm.OrderBy("created_at ASC"),
				qm.Limit(sizeInt),
			)
			all, err2 = result.All(c.Request.Context(), db)
			if err2 != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err2.Error(), nil)
				return
			}
		} else {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid filter parameter", nil)
			return
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(all) > 0 {
			lastCursor = all[len(all)-1].SongInfoID
		}

		// PlaylistAddResponseList 초기화
		playlistAddResponseList := make([]PlaylistAddResponse, 0)

		// all을 순회하며 필요한 정보 추출
		for _, v := range all {
			tempSong := mysql.SongInfos(qm.Where("song_info_id = ?", v.SongInfoID))
			row, err := tempSong.One(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			response := PlaylistAddResponse{
				SongName:   row.SongName,
				SingerName: row.ArtistName,
				SongNumber: row.SongNumber,
				SongInfoId: row.SongInfoID,
				Album:      row.Album.String,
				IsMr:       row.IsMR.Bool,
				IsLive:     row.IsLive.Bool,
				MelonLink:  CreateMelonLinkByMelonSongId(row.MelonSongID),
			}
			playlistAddResponseList = append(playlistAddResponseList, response)
		}

		getPlayListV2Response := GetPlayListV2Response{
			PlayListResponse: playlistAddResponseList,
			LastCursor:       lastCursor,
		}

		// 성공적으로 처리된 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "success", getPlayListV2Response)
	}
}
