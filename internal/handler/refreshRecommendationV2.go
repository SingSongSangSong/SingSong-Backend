package handler

import (
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"net/http"
)

type refreshRequestV2 struct {
	Tag  string `json:"tag"`
	Page int    `json:"page"`
}

type refreshResponseV2 struct {
	Song     []refreshResponse `json:"songs"`
	NextPage int               `json:"nextPage"`
}

// RefreshRecommendationV2 godoc
// @Summary      새로고침 노래 추천V2
// @Description  태그에 해당하는 노래 목록을 보여줍니다. 첫페이지는 1입니당!
// @Tags         Recommendation
// @Accept       json
// @Produce      json
// @Param        songs   body      refreshRequestV2  true  "태그"
// @Success      200 {object} pkg.BaseResponseStruct{data=refreshResponseV2} "성공"
// @Router       /v2/recommend/refresh [post]
// @Security BearerAuth
func RefreshRecommendationV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		request := &refreshRequestV2{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		column, err := MapTagToColumnV4(request.Tag)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid tag", nil)
			return
		}

		page := request.Page
		if page < 1 {
			page = 1 // 최소 1페이지부터 시작
		}
		offset := (page - 1) * pageSize

		//todo: tj_score, 댓글, keep, 멜론 좋아요 종합적으로 반영
		query := fmt.Sprintf(`
			SELECT * FROM (
				SELECT 
					si.song_info_id, si.song_number, si.song_name, si.artist_name, 
					si.album, si.is_mr, si.is_live, si.melon_song_id, si.lyrics_video_link, si.tj_youtube_link, si.tj_score,
					COUNT(DISTINCT c.comment_id) AS comment_count,
					COUNT(DISTINCT ks.keep_song_id) AS keep_count,
					EXISTS (
						SELECT 1 FROM keep_song WHERE song_info_id = si.song_info_id AND keep_list_id IN (
							SELECT keep_list_id FROM keep_list WHERE member_id = ? AND deleted_at IS NULL
						)
					) AS is_keep
				FROM song_info si
				LEFT JOIN comment c ON si.song_info_id = c.song_info_id AND c.deleted_at IS NULL
				LEFT JOIN keep_song ks ON si.song_info_id = ks.song_info_id AND ks.deleted_at IS NULL
				WHERE %s = TRUE
				GROUP BY si.song_info_id
			) AS result
			ORDER BY (result.tj_score + result.keep_count + result.comment_count) DESC, result.song_info_id DESC
			LIMIT ? OFFSET ?
		`, column)

		rows, err := db.Query(query, memberId, pageSize, offset)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close()

		refreshSongs := make([]refreshResponse, 0, pageSize)
		for rows.Next() {
			var songInfoId int64
			var songNumber int
			var songName string
			var artistName string
			var album null.String
			var isMr null.Bool
			var isLive null.Bool
			var melonSongId null.String
			var commentCount int
			var keepCount int
			var isKeep bool
			var lyricsLink null.String
			var tjLink null.String
			var tj_score int

			err := rows.Scan(
				&songInfoId, &songNumber, &songName, &artistName,
				&album, &isMr, &isLive, &melonSongId,
				&lyricsLink, &tjLink, &tj_score,
				&commentCount, &keepCount, &isKeep,
			)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			refreshSongs = append(refreshSongs, refreshResponse{
				SongNumber:        songNumber,
				SongName:          songName,
				SingerName:        artistName,
				Album:             album.String,
				IsKeep:            isKeep,
				SongInfoId:        songInfoId,
				IsMr:              isMr.Bool,
				IsLive:            isLive.Bool,
				KeepCount:         keepCount,
				CommentCount:      commentCount,
				MelonLink:         CreateMelonLinkByMelonSongId(melonSongId),
				LyricsYoutubeLink: lyricsLink.String,
				TJYoutubeLink:     tjLink.String,
			})
		}

		pkg.BaseResponse(c, http.StatusOK, "success", refreshResponseV2{Song: refreshSongs, NextPage: page + 1})
	}
}
