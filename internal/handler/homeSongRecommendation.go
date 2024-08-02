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

type homeSongResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
	SongInfoId int64    `json:"songId"`
	Album      string   `json:"album"`
}

var (
	songInfoIds = []int64{4166, 8525, 46872, 57127, 46375}
)

// HomeSongRecommendation godoc
// @Summary      노래 추천 5곡
// @Description  앨범 이미지와 함께 노래를 추천
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]homeSongResponse} "성공"
// @Router       /recommend/home/songs [get]
func HomeSongRecommendation(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoIdsInterface := make([]interface{}, len(songInfoIds))
		for i, id := range songInfoIds {
			songInfoIdsInterface[i] = id
		}

		all, err := mysql.SongInfos(qm.WhereIn("song_info_id in ?", songInfoIdsInterface...)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		songs := make([]homeSongResponse, 0, len(all))
		for _, s := range all {
			tags := strings.Split(s.Tags.String, ",")
			for i := range tags {
				tags[i] = strings.TrimSpace(tags[i])
			}
			korean, err := MapTagsEnglishToKorean(tags)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			songs = append(songs, homeSongResponse{
				SongNumber: s.SongNumber,
				SongName:   s.SongName,
				SingerName: s.ArtistName,
				Tags:       korean,
				SongInfoId: s.SongInfoID,
				Album:      s.Album.String,
			})
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", songs)
	}
}
