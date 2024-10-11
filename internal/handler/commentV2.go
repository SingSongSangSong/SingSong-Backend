package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"time"
)

// DeleteComment godoc
// @Summary      해당하는 댓글 삭제하기
// @Description  해당하는 댓글 삭제하기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        commentId   path  int  true  "Comment ID"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v2/comment/{commentId} [delete]
// @Security BearerAuth
func DeleteComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// commentId 가져오기
		commentIdParam := c.Param("commentId")
		commentId, err := strconv.ParseInt(commentIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}

		// Delete member
		_, err = mysql.Comments(qm.Where("comment_id = ? AND member_id = ? AND deleted_at is null", commentId, memberId)).
			UpdateAll(c.Request.Context(), db, mysql.M{
				"deleted_at": time.Now(),
			})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}
