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

type CommentWithRecommentsResponse struct {
	CommentId       int64                           `json:"commentId"`
	Content         string                          `json:"content"`
	IsRecomment     bool                            `json:"isRecomment"`
	ParentCommentId int64                           `json:"parentCommentId"`
	SongInfoId      int64                           `json:"songId"`
	MemberId        int64                           `json:"memberId"`
	Nickname        string                          `json:"nickname"`
	CreatedAt       time.Time                       `json:"createdAt"`
	Likes           int                             `json:"likes"`
	IsLiked         bool                            `json:"isLiked"`
	RecommentsCount int                             `json:"recommentsCount"`
	Recomments      []CommentWithRecommentsResponse `json:"recomments"`
}

type CommentPageV3Response struct {
	CommentCount int64                           `json:"commentCount"`
	Comments     []CommentWithRecommentsResponse `json:"comments"`
	LastCursor   int64                           `json:"lastCursor"`
}

// GetCommentsOnSongV3 godoc
// @Summary      특정 노래의 댓글 목록 가져오기V3(최신순, 오래된순 커서페이징 적용)
// @Description  특정 노래의 댓글 목록 가져오기V3(최신순, 오래된순 커서페이징 적용) - query param이 없으면 디폴트는 최신순 입니다.
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        filter query string false "정렬 기준. 최신순=recent, 오래된순(디폴트)=old"
// @Param        size query string false "한번에 조회할 댓글의 개수. 디폴트값은 10 + @(대댓글수)"
// @Param        cursor query string false "마지막에 조회했던 커서의 commentId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 정렬기준의 가장 처음 댓글부터 줌"
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct{data=CommentPageV3Response} "성공"
// @Router       /v3/songs/{songId}/comments [get]
// @Security BearerAuth
func GetCommentsOnSongV3(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		songIdParam := c.Param("songId")
		songId, err := strconv.Atoi(songIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid songId", nil)
			return
		}

		filter := c.DefaultQuery("filter", "old")
		if filter == "" || (filter != "recent" && filter != "old") {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid filter in query", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", "10")
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "0")
		if filter == "recent" {
			cursorStr = c.DefaultQuery("cursor", "9223372036854775807")
		}
		cursorInt, err := strconv.ParseInt(cursorStr, 10, 64)
		if err != nil || cursorInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		// 댓글 가져오기 (최신순/오래된순)
		orderBy := "comment.comment_id ASC"
		cursorCondition := "comment.comment_id > ?" //기본은 오래된순
		if filter == "recent" {
			orderBy = "comment.comment_id DESC"
			cursorCondition = "comment.comment_id < ?" // 최신순
		}

		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			// 댓글 10개만 가져오기!
			qm.Where("comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0", songId),
			qm.WhereNotIn("comment.member_id NOT IN ?", blockedMemberIds...),
			qm.And(cursorCondition, cursorInt),
			qm.OrderBy(orderBy),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)

		commentCount, err := mysql.Comments(
			qm.Where("comment.song_info_id = ? AND comment.deleted_at IS NULL", songId),
		).Count(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(comments) == 0 {
			var lastCursor int64 = 0
			if filter == "old" {
				lastCursor = cursorInt
			}
			pkg.BaseResponse(c, http.StatusOK, "success", CommentPageResponse{
				commentCount,
				[]CommentWithRecommentsCountResponse{},
				lastCursor,
			})
			return
		}

		// comment_id들만 추출
		commentIDs := make([]interface{}, len(comments))
		for i, comment := range comments {
			commentIDs[i] = comment.CommentID
		}

		// 모든 댓글의 RecommentsCount를 한 번에 조회
		recomments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.Where("comment.deleted_at IS NULL"),
			qm.WhereIn("parent_comment_id IN ?", commentIDs...),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...),
			qm.OrderBy("comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// comment_id들만 추출
		reCommentIDs := make([]interface{}, len(recomments))
		for i, recomment := range recomments {
			reCommentIDs[i] = recomment.CommentID
		}

		searchCommentIDs := append(commentIDs, reCommentIDs...)

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.CommentLikes(
			qm.WhereIn("comment_id IN ?", searchCommentIDs...),
			qm.And("member_id = ?", blockerId),
			qm.And("deleted_at is null"),
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
		recommentsCountMap := make(map[int64]int)
		for _, recomment := range recomments {
			if recomment.ParentCommentID.Valid {
				recommentsCountMap[recomment.ParentCommentID.Int64]++
			}
		}
		recommentsMap := make(map[int64][]CommentWithRecommentsResponse)
		for _, recomment := range recomments {
			recommentsMap[recomment.ParentCommentID.Int64] = append(recommentsMap[recomment.ParentCommentID.Int64], CommentWithRecommentsResponse{
				CommentId:       recomment.CommentID,
				Content:         recomment.Content.String,
				IsRecomment:     recomment.IsRecomment.Bool,
				ParentCommentId: recomment.ParentCommentID.Int64,
				SongInfoId:      recomment.SongInfoID,
				MemberId:        recomment.MemberID,
				Nickname:        recomment.R.Member.Nickname.String,
				CreatedAt:       recomment.CreatedAt.Time,
				Likes:           recomment.Likes.Int,
				IsLiked:         likedCommentMap[recomment.CommentID],
			})
		}

		// Initialize a slice to hold all comments
		var topLevelComments []CommentWithRecommentsResponse
		// Add all top-level comments (those without parent comments) to the slice
		for _, comment := range comments {
			topLevelComments = append(topLevelComments, CommentWithRecommentsResponse{
				CommentId:       comment.CommentID,
				Content:         comment.Content.String,
				IsRecomment:     comment.IsRecomment.Bool,
				ParentCommentId: comment.ParentCommentID.Int64,
				SongInfoId:      comment.SongInfoID,
				MemberId:        comment.MemberID,
				Nickname:        comment.R.Member.Nickname.String,
				CreatedAt:       comment.CreatedAt.Time,
				Likes:           comment.Likes.Int,
				IsLiked:         likedCommentMap[comment.CommentID],
				RecommentsCount: recommentsCountMap[comment.CommentID],
				Recomments:      recommentsMap[comment.CommentID],
			})
		}

		response := CommentPageV3Response{
			CommentCount: commentCount,
			Comments:     topLevelComments,
			LastCursor:   topLevelComments[len(topLevelComments)-1].CommentId,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", response)
	}
}