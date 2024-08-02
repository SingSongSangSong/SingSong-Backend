package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

type songReviewOptionGetResponse struct {
	SongReviewOptionId int64  `json:"songReviewOptionId"`
	Title              string `json:"title"`
	Count              int    `json:"count"`
	Selected           bool   `json:"selected"`
}

// GetSongReview godoc
// @Summary      노래 평가를 조회합니다.
// @Description  노래 평가를 조회합니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songNumber path string true "노래 번호"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]songReviewOptionGetResponse} "성공"
// @Router       /songs/{songNumber}/reviews [get]
// @Security BearerAuth
func GetSongReview(db *sql.DB) gin.HandlerFunc {
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

		one, err := mysql.SongInfos(qm.Where("song_number = ?", songNumber)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no song", nil)
			return
		}

		all, err := mysql.SongReviews(qm.Where("song_info_id = ?", one.SongInfoID)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		options, err := mysql.SongReviewOptions().All(c, db)

		response := make([]songReviewOptionGetResponse, 0, len(options))
		for _, option := range options {
			response = append(response, songReviewOptionGetResponse{
				SongReviewOptionId: option.SongReviewOptionID,
				Title:              option.Title.String,
				Count:              0,
				Selected:           false,
			})
		}
		if len(all) != 0 {
			for _, review := range all {
				for i, option := range response {
					if review.SongReviewOptionID == option.SongReviewOptionId {
						response[i].Count++
						if review.MemberID == memberId {
							response[i].Selected = true
						}
						continue
					}
				}
			}
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type songReviewOptionPutRequest struct {
	SongReviewOptionId int64 `json:"songReviewOptionId"`
}

// PutSongReview godoc
// @Summary      노래 평가를 등록/수정합니다.
// @Description  노래 평가를 등록/수정합니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songNumber path string true "노래 번호"
// @Param		 songReview body songReviewOptionPutRequest true "songReviewOptionId"
// @Success      200 "성공"
// @Router       /songs/{songNumber}/reviews [put]
// @Security BearerAuth
func PutSongReview(db *sql.DB) gin.HandlerFunc {
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

		value2, exists := c.Get("birthyear")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - birthyear not found", nil)
			return
		}

		birthyear, ok := value2.(null.Int)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - birthyear not type null.int", nil)
			return
		}

		value3, exists := c.Get("gender")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not found", nil)
			return
		}

		gender, ok := value3.(null.String)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not type null.String", nil)
			return
		}

		var request songReviewOptionPutRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}
		one, err := mysql.SongInfos(qm.Where("song_number = ?", songNumber)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// db에 songReview가 있다면 update, 없다면 insert
		exist, err := mysql.SongReviews(qm.Where("member_id = ?", memberId), qm.And("song_info_id = ?", one.SongInfoID)).Exists(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		if exist {
			_, err := mysql.SongReviews(qm.Where("member_id = ?", memberId), qm.And("song_info_id = ?", one.SongInfoID)).UpdateAll(c, db, mysql.M{"song_review_option_id": request.SongReviewOptionId})
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		} else {
			review := mysql.SongReview{
				SongInfoID:         one.SongInfoID,
				MemberID:           memberId,
				SongReviewOptionID: request.SongReviewOptionId,
				Gender:             gender,
				Birthyear:          birthyear,
			}

			if err := review.Insert(c, db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}

// DeleteSongReview godoc
// @Summary      노래 평가를 삭제합니다.
// @Description  노래 평가를 삭제합니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songNumber path string true "노래 번호"
// @Success      200 "성공"
// @Router       /songs/{songNumber}/reviews [delete]
// @Security BearerAuth
func DeleteSongReview(db *sql.DB) gin.HandlerFunc {
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

		one, err := mysql.SongInfos(qm.Where("song_number = ?", songNumber)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		_, err = mysql.SongReviews(qm.Where("member_id = ?", memberId), qm.And("song_info_id = ?", one.SongInfoID)).DeleteAll(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}