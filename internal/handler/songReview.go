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
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]songReviewOptionGetResponse} "성공"
// @Router       /v1/songs/{songId}/reviews [get]
// @Security BearerAuth
func GetSongReview(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
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

		all, err := mysql.SongReviews(qm.Where("song_info_id = ?", songInfoId), qm.And("deleted_at IS NULL")).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		options, err := mysql.SongReviewOptions().All(c.Request.Context(), db)

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
// @Param        songId path string true "songId"
// @Param		 songReview body songReviewOptionPutRequest true "songReviewOptionId"
// @Success      200 "성공"
// @Router       /v1/songs/{songId}/reviews [put]
// @Security BearerAuth
func PutSongReview(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
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

		value2, exists := c.Get("birthYear")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - birthYear not found", nil)
			return
		}
		birthYearStr := value2.(string)
		birthYear, err := strconv.Atoi(birthYearStr)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not found", nil)
			return
		}

		value3, exists := c.Get("gender")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - gender not found", nil)
			return
		}
		gender := value3.(string)

		var request songReviewOptionPutRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}
		one, err := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId)).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// soft delete
		_, err = mysql.SongReviews(
			qm.Where("member_id = ?", memberId), qm.And("song_info_id = ?", one.SongInfoID), qm.And("deleted_at IS NULL"),
		).UpdateAll(c.Request.Context(), db, mysql.M{"deleted_at": null.TimeFrom(time.Now())})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// insert
		review := mysql.SongReview{
			SongInfoID:         one.SongInfoID,
			MemberID:           memberId,
			SongReviewOptionID: request.SongReviewOptionId,
			Gender:             null.StringFrom(gender),
			Birthyear:          null.IntFrom(birthYear),
		}

		if err := review.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		go func(db *sql.DB, memberId interface{}, songReviewOptionId int64) {
			ctx := context.Background()
			option, err2 := mysql.SongReviewOptions(qm.Where("song_review_option_id = ?", songReviewOptionId)).One(ctx, db)
			if err2 != nil {
				log.Printf("failed to get song review option: " + err2.Error())
				return
			}
			logMemberAction(db, memberId, "REVIEW_"+option.Enum.String, 3, songInfoId)
		}(db, value, request.SongReviewOptionId)

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}

// DeleteSongReview godoc
// @Summary      노래 평가를 삭제합니다.
// @Description  노래 평가를 삭제합니다.
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songId path string true "songId"
// @Success      200 "성공"
// @Router       /v1/songs/{songId}/reviews [delete]
// @Security BearerAuth
func DeleteSongReview(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
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

		// soft delete
		_, err := mysql.SongReviews(
			qm.Where("member_id = ?", memberId), qm.And("song_info_id = ?", songInfoId), qm.And("deleted_at IS NULL"),
		).UpdateAll(c.Request.Context(), db, mysql.M{"deleted_at": null.TimeFrom(time.Now())})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}
