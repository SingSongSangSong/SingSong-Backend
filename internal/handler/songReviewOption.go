package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"net/http"
)

type songReviewOptionAddRequest struct {
	Title string `json:"title"`
}

// AddSongReviewOption godoc
// @Summary      노래 평가 선택지를 추가합니다.
// @Description  노래 평가 선택지를 추가합니다.
// @Tags         Song review option CR for admin
// @Accept       json
// @Produce      json
// @Param        songReviewOptionAddRequest body songReviewOptionAddRequest true "평가 선택지"
// @Success      200 "성공"
// @Router       /song-review-options [post]
func AddSongReviewOption(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &songReviewOptionAddRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// Convert the string to null.String
		title := null.StringFrom(request.Title)
		option := mysql.SongReviewOption{Title: title}

		if err := option.Insert(c, db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}

// ListSongReviewOptions godoc
// @Summary      노래 평가 선택지를 모두 조회합니다.
// @Description  노래 평가 선택지를 모두 조회합니다.
// @Tags         Song review option CR for admin
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]string} "성공"
// @Router       /song-review-options [get]
func ListSongReviewOptions(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		all, err := mysql.SongReviewOptions().All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		response := make([]string, 0, len(all))
		for _, option := range all {
			response = append(response, option.Title.String)
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
