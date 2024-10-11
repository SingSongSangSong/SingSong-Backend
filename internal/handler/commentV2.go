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

type MyCommentPageResponse struct {
	Comments   []MyComment `json:"comments"`
	LastCursor int64       `json:"lastCursor"`
}

type MyComment struct {
	CommentId       int64           `json:"commentId"`
	Content         string          `json:"content"`
	IsRecomment     bool            `json:"isRecomment"`
	ParentCommentId int64           `json:"parentCommentId"`
	CreatedAt       time.Time       `json:"createdAt"`
	Likes           int             `json:"likes"`
	IsLiked         bool            `json:"isLiked"`
	Song            SongOfMyComment `json:"song"`
}

type SongOfMyComment struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
	IsLive     bool   `json:"isLive"`
	MelonLink  string `json:"melonLink"`
}

// GetMyComments godoc
// @Summary      내가 쓴 댓글 모아보기
// @Description  내가 쓴 댓글 모아보기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        size   query      int  false  "size"
// @Param        cursor   query      int  false  "cursor"
// @Success      200 {object} pkg.BaseResponseStruct{data=MyCommentPageResponse} "성공"
// @Router       /v2/comment/my [get]
// @Security BearerAuth
func GetMyComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		comments, err := mysql.Comments(
			qm.Where("member_id = ?", memberId),
			qm.Where("deleted_at is null"),
			qm.Where("comment_id < ?", cursorInt),
			qm.OrderBy("comment_id DESC"), // 최신 순 정렬
			qm.Limit(sizeInt),             // 최신 size개의 댓글만 가져옴
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// song_info_id들만 추출
		songInfoIDs := make([]interface{}, len(comments))
		for i, comment := range comments {
			songInfoIDs[i] = comment.SongInfoID
		}

		// 노래 조회
		songs, err := mysql.SongInfos(
			qm.WhereIn("song_info_id IN ?", songInfoIDs...),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// song 정보를 를 맵으로 저장
		songMap := make(map[int64]*mysql.SongInfo)
		for _, song := range songs {
			songMap[song.SongInfoID] = song
		}

		// comment_id들만 추출
		commentIDs := make([]interface{}, len(comments))
		for i, comment := range comments {
			commentIDs[i] = comment.CommentID
		}

		// 댓글 좋아요 여부 조회
		likes, err := mysql.CommentLikes(
			qm.WhereIn("comment_id IN ?", commentIDs...),
			qm.And("member_id = ?", memberId),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 좋아요를 누른 comment_id를 맵으로 저장 (빠른 조회를 위해)
		likedCommentMap := make(map[int64]bool)
		for _, like := range likes {
			likedCommentMap[like.CommentID] = true
		}

		myComments := make([]MyComment, 0, sizeInt)

		for _, comment := range comments {
			song := songMap[comment.SongInfoID]
			myComments = append(myComments, MyComment{
				CommentId:       comment.CommentID,
				Content:         comment.Content.String,
				IsRecomment:     comment.IsRecomment.Bool,
				ParentCommentId: comment.ParentCommentID.Int64,
				CreatedAt:       comment.CreatedAt.Time,
				Likes:           comment.Likes.Int,
				IsLiked:         likedCommentMap[comment.CommentID],
				Song: SongOfMyComment{
					song.SongNumber,
					song.SongName,
					song.ArtistName,
					song.SongInfoID,
					song.Album.String,
					song.IsMR.Bool,
					song.IsLive.Bool,
					CreateMelonLinkByMelonSongId(song.MelonSongID),
				},
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(myComments) > 0 {
			lastCursor = myComments[len(myComments)-1].CommentId
		}

		response := MyCommentPageResponse{
			Comments:   myComments,
			LastCursor: lastCursor,
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
