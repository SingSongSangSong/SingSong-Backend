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
	SongInfoId  int64    `json:"songId"`
	Album       string   `json:"Album"`
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
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]songInfoResponse} "성공"
// @Router       /songs/{songId} [get]
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
		one, err := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no song", nil)
			return
		}

		//유저의 keep 여부 조회
		all, err := mysql.KeepLists(qm.Where("member_id = ?", memberId)).All(c, db)
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
		).Exists(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		response := songInfoResponse{
			SongNumber:  one.SongNumber,
			SongName:    one.SongName,
			SingerName:  one.ArtistName,
			Tags:        parseTags(one.Tags.String),
			SongInfoId:  one.SongInfoID,
			Album:       one.Album.String,
			Octave:      one.Octave.String,
			Description: "", //todo:
			IsKeep:      isKeep,
		}

		// 비동기적으로 member_action 저장
		go logMemberAction(db, memberId, "CLICK", 0.5, songInfoId)

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
