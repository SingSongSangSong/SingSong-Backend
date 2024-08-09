package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"strconv"
	"time"
)

type PlaylistAddRequest struct {
	SongInfoIds []int `json:"songId"`
}

type PlaylistAddResponse struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
}

// GoRoutine으로 회원가입시에 플레이리스트를 생성한다 (context따로 가져와야함)
func CreatePlaylist(db *sql.DB, keepName string, memberId int64) {
	// 플레이리스트 생성
	m := mysql.KeepList{MemberID: memberId, KeepName: null.StringFrom(keepName)}
	err := m.Insert(context.Background(), db, boil.Infer())
	if err != nil {
		log.Printf("Error inserting Playlist: %v", err)
	}
}

// 플레이리스트에 노래리스트 추가
// AddSongsToPlaylist godoc
// @Summary      플레이리스트에 노래를 추가한다
// @Description  노래들을 하나씩 플레이리스트에 추가한 후 적용된 플레이리스트의 노래들을 리턴한다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        PlaylistAddRequest  body   PlaylistAddRequest  true  "노래 리스트"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
// @Router       /keep [post]
// @Security BearerAuth
func AddSongsToPlaylist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlistRequest := &PlaylistAddRequest{}
		if err := c.ShouldBindJSON(&playlistRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Playlist 정보 가져오기
		m := mysql.KeepLists(qm.Where("member_id = ?", memberId))
		playlistRow, errors := m.One(c, db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
			return
		}

		// 노래 정보들 가져오기
		for _, songInfoId := range playlistRequest.SongInfoIds {
			m := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId))
			row, errors := m.One(c, db)
			if errors != nil {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
				return
			}

			// 기존에 같은 keepId와 songTempId가 있는지 확인
			existsQuery := mysql.KeepSongs(qm.Where("keep_list_id = ? AND song_info_id = ? AND deleted_at IS NULL", playlistRow.KeepListID, row.SongInfoID))
			existingRow, err := existsQuery.One(c, db)
			if err == nil && existingRow != nil {
				// 이미 존재하면 추가하지 않고 계속 진행
				continue
			}

			keepSong := mysql.KeepSong{KeepListID: playlistRow.KeepListID, SongInfoID: row.SongInfoID, SongNumber: row.SongNumber}
			err = keepSong.Insert(c, db, boil.Infer())
			if err != nil {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
				return
			}
		}

		go func(db *sql.DB, memberId interface{}, songInfoIds []int) {
			songInfoIdsStr := make([]string, len(songInfoIds))
			for i, v := range songInfoIds {
				songInfoIdsStr[i] = strconv.Itoa(v)
			}
			logMemberAction(db, memberId, "KEEP", 2, songInfoIdsStr...)
		}(db, memberId, playlistRequest.SongInfoIds)

		result := mysql.KeepSongs(qm.Where("keep_list_id = ?", playlistRow.KeepListID))
		all, err2 := result.All(c, db)
		if err2 != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err2.Error(), nil)
			return
		}

		PlaylistAddResponseList := make([]PlaylistAddResponse, 0)

		for _, v := range all {
			tempSong := mysql.SongInfos(qm.Where("song_info_id = ?", v.SongInfoID))
			row, errors := tempSong.One(c, db)
			if errors != nil {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
				return
			}
			response := PlaylistAddResponse{SongName: row.SongName, SingerName: row.ArtistName, SongNumber: row.SongNumber, SongInfoId: row.SongInfoID}
			PlaylistAddResponseList = append(PlaylistAddResponseList, response)
		}

		pkg.BaseResponse(c, http.StatusOK, "success", PlaylistAddResponseList)
	}
}

type SongDeleteFromPlaylistRequest struct {
	SongInfoIds []int `json:"songIds"`
}

// 플레이리스트에 노래리스트 삭제
// DeleteSongsFromPlaylist godoc
// @Summary      플레이리스트에 노래를 제거한다
// @Description  노래들을 하나씩 플레이리스트에서 삭제한다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Param        SongDeleteFromPlaylistRequest  body      SongDeleteFromPlaylistRequest  true  "노래 리스트"
// @Success      200 {object} pkg.BaseResponseStruct{data=PlaylistAddResponse} "성공"
// @Router       /keep [delete]
// @Security BearerAuth
func DeleteSongsFromPlaylist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songDeleteFromPlaylistRequest := &SongDeleteFromPlaylistRequest{}
		if err := c.ShouldBindJSON(&songDeleteFromPlaylistRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		memberId, err := c.Get("memberId")
		if err != true {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Playlist정보 가져오기
		m := mysql.KeepLists(qm.Where("member_id = ?", memberId))
		playlistInfo, errors := m.One(c, db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
			return
		}

		// 노래 삭제
		for _, songInfoId := range songDeleteFromPlaylistRequest.SongInfoIds {
			_, err := mysql.KeepSongs(
				qm.Where("keep_list_id = ? AND song_info_id = ? AND deleted_at IS NULL", playlistInfo.KeepListID, songInfoId),
			).UpdateAll(c, db, mysql.M{"deleted_at": null.TimeFrom(time.Now())})
			if err != nil {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			}
		}

		// 응답에 keep 목록 넣기
		all, errors := mysql.KeepSongs(qm.Where("keep_list_id = ? AND deleted_at IS NULL", playlistInfo.KeepListID)).All(c, db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+errors.Error(), nil)
			return
		}

		keepSongs := make([]PlaylistAddResponse, 0)

		for _, v := range all {
			tempSong := mysql.SongInfos(qm.Where("song_info_id = ?", v.SongInfoID))
			row, errors := tempSong.One(c, db)
			if errors != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+errors.Error(), nil)
				return
			}
			response := PlaylistAddResponse{SongName: row.SongName, SingerName: row.ArtistName, SongNumber: row.SongNumber, SongInfoId: row.SongInfoID}
			keepSongs = append(keepSongs, response)
		}
		pkg.BaseResponse(c, http.StatusOK, "success", keepSongs)
	}
}

// 플레이리스트에 노래리스트 조회
// GetSongsFromPlaylist godoc
// @Summary      플레이리스트에 노래를 가져온다
// @Description  플레이리스트에 있는 노래들을 가져온다
// @Tags         Playlist
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
// @Router       /keep [get]
// @Security BearerAuth
func GetSongsFromPlaylist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, err := c.Get("memberId")
		if err != true {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Playlist정보 가져오기
		m := mysql.KeepLists(qm.Where("member_id = ?", memberId))
		playlistInfo, errors := m.One(c, db)
		if errors != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
			return
		}

		result := mysql.KeepSongs(qm.Where("keep_list_id = ? AND deleted_at IS NULL", playlistInfo.KeepListID))
		all, err2 := result.All(c, db)
		if err2 != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err2.Error(), nil)
			return
		}

		PlaylistAddResponseList := make([]PlaylistAddResponse, 0)

		for _, v := range all {
			tempSong := mysql.SongInfos(qm.Where("song_info_id = ?", v.SongInfoID))
			row, errors := tempSong.One(c, db)
			if errors != nil {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - "+errors.Error(), nil)
				return
			}
			response := PlaylistAddResponse{SongName: row.SongName, SingerName: row.ArtistName, SongNumber: row.SongNumber, SongInfoId: row.SongInfoID}
			PlaylistAddResponseList = append(PlaylistAddResponseList, response)
		}

		pkg.BaseResponse(c, http.StatusOK, "success", PlaylistAddResponseList)
	}
}

// 플레이리스트에 노래리스트 수정
