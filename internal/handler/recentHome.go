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

// GetLatestSearchApi 검색화면 최근 검색어
// @Summary      검색화면 최근 검색어 가져오기
// @Description  검색화면 최근 검색어 가져오기. 쿼리 파라미터인 size를 별도로 지정하지 않으면 default size = 10
// @Tags         Recent
// @Accept       json
// @Produce      json
// @Param        size   query      int  false  "size"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]interface{}} "Success"
// @Router       /v1/recent/search [get]
func GetLatestSearchApi(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = "10" //default value
		}

		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot convert size to int", nil)
			return
		}

		// 최근 검색어 가져오기
		latestSearch, err := mysql.SearchLogs(
			qm.InnerJoin(""+
				"(SELECT search_text, MAX(created_at) AS max_created_at "+
				"FROM search_log "+
				"GROUP BY search_text) AS latest_search "+
				"ON search_log.search_text = latest_search.search_text AND search_log.created_at = latest_search.max_created_at"),
			qm.OrderBy("created_at DESC"),
			qm.Limit(size)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// interface
		response := make([]interface{}, len(latestSearch))
		for i, search := range latestSearch {
			response[i] = search.SearchText
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "success", response)
	}
}

// GetRecentKeepSongs
// @Summary      최근 Keep한 노래 목록 가져오기
// @Description  최근 Keep한 노래 목록 가져오기. 쿼리 파라미터인 size를 별도로 지정하지 않으면 default size = 10
// @Tags         Recent
// @Accept       json
// @Produce      json
// @Param        size   query      int  false  "size"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]SongSearchInfoV2Response} "Success"
// @Router       /v1/recent/keep [get]
func GetRecentKeepSongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = "10" //default value
		}

		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot convert size to int", nil)
			return
		}

		// 저장한 노래 가져오기
		likeSongs, err := mysql.KeepSongs(
			qm.InnerJoin(""+
				"(SELECT song_info_id, MAX(created_at) AS max_created_at "+
				"FROM keep_song "+
				"WHERE deleted_at IS NULL "+
				"GROUP BY song_info_id) AS latest_songs "+
				"ON keep_song.song_info_id = latest_songs.song_info_id AND keep_song.created_at = latest_songs.max_created_at"),
			qm.OrderBy("created_at DESC"),
			qm.Limit(size),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// interface
		response := make([]SongSearchInfoV2Response, len(likeSongs))
		songInfoIds := make([]interface{}, len(likeSongs))

		for i, likeSong := range likeSongs {
			songInfoIds[i] = likeSong.SongInfoID
		}

		// 노래 정보 가져오기
		songInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 노래 정보 맵 생성
		for i, songInfo := range songInfos {
			response[i] = SongSearchInfoV2Response{
				SongNumber:        songInfo.SongNumber,
				SongInfoId:        songInfo.SongInfoID,
				SongName:          songInfo.SongName,
				SingerName:        songInfo.ArtistName,
				Album:             songInfo.Album.String,
				IsMr:              songInfo.IsMR.Bool,
				IsLive:            songInfo.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(songInfo.MelonSongID),
				LyricsYoutubeLink: songInfo.LyricsVideoLink.String,
				TJYoutubeLink:     songInfo.TJYoutubeLink.String,
				LyricsVideoID:     ExtractVideoID(songInfo.LyricsVideoLink.String),
				TJVideoID:         ExtractVideoID(songInfo.TJYoutubeLink.String),
			}
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "success", response)
	}
}

// GetRecentCommentsongs
// @Summary      최근 댓글 단 노래 목록 가져오기
// @Description  최근 댓글 단 노래 목록 가져오기. 쿼리 파라미터인 size를 별도로 지정하지 않으면 default size = 10
// @Tags         Recent
// @Accept       json
// @Produce      json
// @Param        size   query      int  false  "size"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]SongSearchInfoV2Response} "Success"
// @Router       /v1/recent/comment [get]
func GetRecentCommentsongs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = "10" //default value
		}

		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot convert size to int", nil)
			return
		}

		// 댓글 단 노래 가져오기
		commentSongs, err := mysql.Comments(
			qm.InnerJoin(""+
				"(SELECT song_info_id, MAX(created_at) AS max_created_at "+
				"FROM comment "+
				"WHERE deleted_at IS NULL "+
				"GROUP BY song_info_id) AS latest_comment "+
				"ON comment.song_info_id = latest_comment.song_info_id AND comment.created_at = latest_comment.max_created_at"),
			qm.OrderBy("created_at DESC"),
			qm.Limit(size)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// interface
		response := make([]SongSearchInfoV2Response, len(commentSongs))
		songInfoIds := make([]interface{}, len(commentSongs))

		for i, commentSong := range commentSongs {
			songInfoIds[i] = commentSong.SongInfoID
		}

		// 노래 정보 가져오기
		songInfos, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 노래 정보 맵 생성
		for i, songInfo := range songInfos {
			response[i] = SongSearchInfoV2Response{
				SongNumber:        songInfo.SongNumber,
				SongInfoId:        songInfo.SongInfoID,
				SongName:          songInfo.SongName,
				SingerName:        songInfo.ArtistName,
				Album:             songInfo.Album.String,
				IsMr:              songInfo.IsMR.Bool,
				IsLive:            songInfo.IsLive.Bool,
				MelonLink:         CreateMelonLinkByMelonSongId(songInfo.MelonSongID),
				LyricsYoutubeLink: songInfo.LyricsVideoLink.String,
				TJYoutubeLink:     songInfo.TJYoutubeLink.String,
				LyricsVideoID:     ExtractVideoID(songInfo.LyricsVideoLink.String),
				TJVideoID:         ExtractVideoID(songInfo.TJYoutubeLink.String),
			}
		}

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "success", response)
	}
}
