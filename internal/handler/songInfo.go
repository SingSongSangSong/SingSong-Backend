package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strings"
)

type songInfoResponse struct {
	SongNumber        int      `json:"songNumber"`
	SongName          string   `json:"songName"`
	SingerName        string   `json:"singerName"`
	Tags              []string `json:"tags"`
	SongInfoId        int64    `json:"songId"`
	Album             string   `json:"album"`
	Octave            string   `json:"octave"`
	Description       string   `json:"description"`
	IsKeep            bool     `json:"isKeep"`
	KeepCount         int64    `json:"keepCount"`
	CommentCount      int64    `json:"commentCount"`
	IsMr              bool     `json:"isMr"`
	IsLive            bool     `json:"isLive"`
	MelonLink         string   `json:"melonLink"`
	LyricsYoutubeLink string   `json:"lyricsYoutubeLink"`
	TJYoutubeLink     string   `json:"tjYoutubeLink"`
	LyricsVideoID     string   `json:"lyricsVideoId"`
	TJVideoID         string   `json:"tjVideoId"`
}

// GetSongInfo godoc
// @Summary      노래 상세 정보를 조회합니다
// @Description  노래 상세 정보를 조회합니다
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]songInfoResponse} "성공"
// @Router       /v1/songs/{songId} [get]
// @Security BearerAuth
func GetSongInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
			return
		}

		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		//노래 정보 조회
		one, err := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId)).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no song", nil)
			return
		}

		//유저의 keep 여부 조회
		all, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepListIds := make([]interface{}, len(all))
		for i, keep := range all {
			keepListIds[i] = keep.KeepListID
		}
		isKeep, err := mysql.KeepSongs(
			qm.WhereIn("keep_list_id in ?", keepListIds...),
			qm.And("song_info_id = ?", one.SongInfoID),
			qm.And("deleted_at IS NULL"),
		).Exists(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		commentCount, err := mysql.Comments(qm.Where("song_info_id = ? AND deleted_at is null", one.SongInfoID)).Count(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		keepCount, err := mysql.KeepSongs(qm.Where("song_info_id = ? AND deleted_at is null", one.SongInfoID)).Count(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		response := songInfoResponse{
			SongNumber:        one.SongNumber,
			SongName:          one.SongName,
			SingerName:        one.ArtistName,
			SongInfoId:        one.SongInfoID,
			Tags:              []string{}, //todo: tags
			Album:             one.Album.String,
			Octave:            one.Octave.String,
			Description:       "", //todo:
			IsKeep:            isKeep,
			CommentCount:      commentCount,
			KeepCount:         keepCount,
			IsMr:              one.IsMR.Bool,
			IsLive:            one.IsLive.Bool,
			MelonLink:         CreateMelonLinkByMelonSongId(one.MelonSongID),
			LyricsYoutubeLink: one.LyricsVideoLink.String,
			TJYoutubeLink:     one.TJYoutubeLink.String,
			LyricsVideoID:     ExtractVideoID(one.LyricsVideoLink.String),
			TJVideoID:         ExtractVideoID(one.TJYoutubeLink.String),
		}

		// 비동기적으로 member_action 저장
		go logMemberAction(db, memberId, "CLICK", 0.5, songInfoId)

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

func CreateMelonLinkByMelonSongId(melonSongId null.String) string {
	if melonSongId.Valid && melonSongId.String != "" {
		return "https://www.melon.com/song/detail.htm?songId=" + melonSongId.String
	}
	return "https://www.melon.com/"
}

func parseTags(tagString string) []string {
	tags := make([]string, 0)
	if trimmedTags := strings.TrimSpace(tagString); trimmedTags != "" {
		for _, tag := range strings.Split(trimmedTags, ",") {
			trimmedTag := strings.TrimSpace(tag)
			if trimmedTag != "" {
				tags = append(tags, trimmedTag)
			}
		}
	}
	return tags
}

// GetLinkBySongId godoc
// @Summary      songId로 link를 조회합니다.
// @Description  songId로 link를 조회합니다.
// @Tags         Link
// @Accept       json
// @Produce      json
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct(data=string) "성공"
// @Router       /v1/songs/{songId}/link [get]
func GetLinkBySongInfoId(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
			return
		}

		link := GetMelonLink(c, songInfoId, db)

		pkg.BaseResponse(c, http.StatusOK, "ok", link)
	}
}

func GetMelonLink(c *gin.Context, songInfoId string, db *sql.DB) string {
	link := "https://www.melon.com/"
	info := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId))
	one, err := info.One(c.Request.Context(), db)
	if err == nil && one.MelonSongID.Valid == true {
		link = "https://www.melon.com/song/detail.htm?songId=" + one.MelonSongID.String
	}
	return link
}
