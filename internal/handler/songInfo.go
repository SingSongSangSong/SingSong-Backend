package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strings"
)

type songInfoResponse struct {
	SongNumber  int      `json:"songNumber"`
	SongName    string   `json:"songName"`
	SingerName  string   `json:"singerName"`
	Tags        []string `json:"tags"`
	SongTempId  int64    `json:"songId"`
	Album       string   `json:"album"`
	Octave      string   `json:"octave"`
	Description string   `json:"description"`
	IsKeep      bool     `json:"isKeep"`
}

// GetSongInfo godoc
// @Summary      노래 상세 정보를 조회합니다
// @Description  노래 상세 정보를 조회합니다
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songNumber path string true "노래 번호"
// @Success      200 {object} pkg.BaseResponseStruct(data=[]songInfoResponse) "성공"
// @Router       /songs/{songNumber} [get]
// @Security BearerAuth
func GetSongInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songNumber := c.Param("songNumber")
		if songNumber == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songNumber in path variable", nil)
			return
		}

		value, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		memberId, ok := value.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not type int64", nil)
			return
		}

		//노래 정보 조회
		one, err := mysql.SongTempInfos(qm.Where("songNumber = ?", songNumber)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no song", nil)
			return
		}

		//유저의 keep 여부 조회
		all, err := mysql.KeepLists(qm.Where("memberId = ?", memberId)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepIds := make([]interface{}, len(all))
		for i, keep := range all {
			keepIds[i] = keep.KeepId
		}
		isKeep, err := mysql.KeepSongs(qm.WhereIn("keepId in ?", keepIds...)).Exists(c, db) // todo:
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		response := songInfoResponse{
			SongNumber:  one.SongNumber,
			SongName:    one.SongName,
			SingerName:  one.ArtistName,
			Tags:        parseTags(one.Tags.String),
			SongTempId:  one.SongTempId,
			Album:       one.Album.String,
			Octave:      one.Octave.String,
			Description: "20대 남성이 가장 많이 부른 노래 Top 1", //todo: 하드 코딩 제거
			IsKeep:      isKeep,
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
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
