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
	"strconv"
	"time"
)

// LikePost godoc
// @Summary      게시글 좋아요 추가
// @Description  게시글 좋아요 추가
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId path string true "postId"
// @Success      200 "성공"
// @Failure      400 "postId param이 잘못 들어왔다면 400 실패"
// @Failure      401 "토큰 인증에 실패했다면 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/{postId}/likes [post]
// @Security BearerAuth
func AddPostLike(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIdStr := c.Param("postId")
		if postIdStr == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find postId in path variable", nil)
			return
		}

		postId, err := strconv.ParseInt(postIdStr, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postId", nil)
			return
		}

		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// 이미 좋아요를 눌렀다면 에러
		exists, err = mysql.PostLikes(
			qm.Where("member_id = ? AND post_id = ? AND deleted_at IS NULL", memberId.(int64), postId),
		).Exists(c.Request.Context(), db)
		if exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - already liked", nil)
			return
		}
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		like := mysql.PostLike{
			MemberID: memberId.(int64),
			PostID:   postId,
		}

		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		//todo post의 likes도 1 증가시켜야 함.

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

// LikePost godoc
// @Summary      게시글 좋아요 해제
// @Description  게시글 좋아요 해제
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId path string true "postId"
// @Success      200 "성공"
// @Failure      400 "postId param이 잘못 들어왔다면 400 실패"
// @Failure      401 "토큰 인증에 실패했다면 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/{postId}/likes [delete]
// @Security BearerAuth
func DeletePostLike(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postId := c.Param("postId")
		if postId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find postId in path variable", nil)
			return
		}

		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		_, err := mysql.PostLikes(
			qm.Where("member_id = ?", memberId),
			qm.And("post_id = ?", postId),
			qm.And("deleted_at IS NULL"),
		).UpdateAll(c.Request.Context(), db, mysql.M{"deleted_at": null.TimeFrom(time.Now())})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		//todo: post의 like도 하나 감소 시켜야 함

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}
