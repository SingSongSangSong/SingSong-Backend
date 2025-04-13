package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"SingSong-Server/internal/repository"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type CommentWithRecommentsResponse struct {
	CommentId       int64                           `json:"commentId"`
	Content         string                          `json:"content"`
	IsRecomment     bool                            `json:"isRecomment"`
	ParentCommentId int64                           `json:"parentCommentId"`
	SongInfoId      int64                           `json:"songId"`
	MemberId        int64                           `json:"memberId"`
	IsWriter        bool                            `json:"isWriter"`
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

		blockedMemberIds, err := repository.GetBlockedMemberIDs(c.Request.Context(), db, blockerId)

		// 댓글 가져오기 (최신순/오래된순)
		orderBy := "comment.comment_id ASC"
		cursorCondition := "comment.comment_id > ?" //기본은 오래된순
		if filter == "recent" {
			orderBy = "comment.comment_id DESC"
			cursorCondition = "comment.comment_id < ?" // 최신순
		}

		var (
			comments     mysql.CommentSlice
			commentCount int64
			err1, err2   error
		)

		wg := sync.WaitGroup{}
		wg.Add(2)

		// 병렬 처리 #1: 댓글 목록과 댓글 개수 동시에
		go func() {
			defer wg.Done()
			comments, err1 = repository.GetTopLevelComments(c, db, songId, cursorCondition, cursorInt, orderBy, sizeInt, blockedMemberIds)
		}()

		go func() {
			defer wg.Done()
			commentCount, err2 = repository.GetCommentCount(c, db, songId)
		}()

		wg.Wait()
		if err1 != nil {
			pkg.SendToSentryWithStack(c, err1)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err1.Error(), nil)
			return
		}
		if err2 != nil {
			pkg.SendToSentryWithStack(c, err2)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err2.Error(), nil)
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
		recomments, err := repository.GetRecomments(c, db, commentIDs, blockedMemberIds)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
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
		likes, err := repository.GetLikesForComments(c, db, searchCommentIDs, blockerId)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
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
				IsWriter:        recomment.MemberID == blockerId,
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
				IsWriter:        comment.MemberID == blockerId,
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
