package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type PostCommentRequest struct {
	ParentCommentId int64  `json:"parentCommentId"`
	PostId          int64  `json:"postId"`
	Content         string `json:"content"`
	IsRecomment     bool   `json:"isRecomment"`
}

type PostCommentResponse struct {
	CommentId       int64                 `json:"commentId"`
	Content         string                `json:"content"`
	IsRecomment     bool                  `json:"isRecomment"`
	ParentCommentId int64                 `json:"parentCommentId"`
	PostId          int64                 `json:"postId"`
	MemberId        int64                 `json:"memberId"`
	Nickname        string                `json:"nickname"`
	CreatedAt       time.Time             `json:"createdAt"`
	Likes           int                   `json:"likes"`
	IsLiked         bool                  `json:"isLiked"`
	PostRecomments  []PostCommentResponse `json:"postRecomments"`
}

// CommentOnPost godoc
// @Summary      PostId에 댓글 달기
// @Description  PostId에 댓글 달기
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        PostCommentRequest   body      PostCommentRequest  true  "postCommentRequest"
// @Success      200 {object} pkg.BaseResponseStruct{data=PostCommentResponse} "성공"
// @Router       /v1/posts/comments [post]
// @Security BearerAuth
func CommentOnPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CommentRequest 받기
		commentRequest := &PostCommentRequest{}
		if err := c.ShouldBindJSON(commentRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// memberId가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		member, err := mysql.Members(qm.Where("member_id = ?", memberId.(int64))).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 댓글 달기
		nulIsRecomment := null.BoolFrom(commentRequest.IsRecomment)
		nullParentCommentId := null.Int64From(commentRequest.ParentCommentId)
		nullContent := null.StringFrom(commentRequest.Content)
		m := mysql.PostComment{MemberID: memberId.(int64), ParentPostCommentID: nullParentCommentId, PostID: commentRequest.PostId, IsRecomment: nulIsRecomment, Content: nullContent, Likes: 0}
		err = m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		commentResponse := PostCommentResponse{
			CommentId:       m.PostCommentID,
			ParentCommentId: m.ParentPostCommentID.Int64,
			PostId:          m.PostID,
			Content:         m.Content.String,
			IsRecomment:     m.IsRecomment.Bool,
			MemberId:        m.MemberID,
			Nickname:        member.Nickname.String,
			CreatedAt:       m.CreatedAt.Time,
			Likes:           m.Likes,
			PostRecomments:  []PostCommentResponse{},
		}

		// 댓글 달기 성공시 댓글 정보 반환
		pkg.BaseResponse(c, http.StatusOK, "success", commentResponse)
	}
}

type GetPostCommentResponse struct {
	TotalPostCommentCount int                   `json:"totalPostCommentCount"`
	PostComments          []PostCommentResponse `json:"postComments"`
	NextPage              int                   `json:"nextPage"`
}

// GetCommentOnPost godoc
// @Summary      Retrieve comments for the specified postId
// @Description  Get comments for a specific post identified by postId with optional page and size query parameters
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId   path     int  true  "Post ID"
// @Param        page query int false "현재 조회할 게시글 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetPostCommentResponse} "Success"
// @Router       /v1/posts/{postId}/comments [get]
// @Security BearerAuth
func GetCommentOnPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve songId from path parameter
		postIdParam := c.Param("postId")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postCommentId", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		pageStr := c.DefaultQuery("page", defaultPage)
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		// OFFSET 계산
		offset := (pageInt - 1) * sizeInt

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		//차단 유저 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		//blocked_member_id 리스트 만들기
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		// Retrieve comments for the specified songId
		postComments, err := mysql.PostComments(
			qm.Load(mysql.PostCommentRels.Member),
			qm.LeftOuterJoin("member ON member.member_id = post_comment.member_id"),
			qm.Where("post_comment.post_id = ? AND post_comment.deleted_at IS NULL", postId),
			qm.WhereNotIn("post_comment.member_id NOT IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.Limit(sizeInt), // limit은 가져올 댓글 수
			qm.Offset(offset), // offset은 몇 번째부터 시작할지
			qm.OrderBy("post_comment.created_at DESC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 전체 댓글 수를 가져오는 쿼리
		totalCommentsCount, err := mysql.PostComments(
			qm.Where("post_comment.post_id = ? AND post_comment.is_recomment = false AND post_comment.deleted_at IS NULL", postId),
			qm.WhereNotIn("post_comment.member_id NOT IN ?", blockedMemberIds...), // 블랙리스트 제외
		).Count(c.Request.Context(), db)

		// 결과가 없는 경우 빈 리스트 반환
		if len(postComments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentResponse{TotalPostCommentCount: int(totalCommentsCount), PostComments: []PostCommentResponse{}, NextPage: pageInt})
			return
		}

		// comment_id들만 추출
		postCommentIDs := make([]interface{}, len(postComments))
		for i, postComment := range postComments {
			postCommentIDs[i] = postComment.PostCommentID
		}

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.PostCommentLikes(
			qm.WhereIn("post_comment_id IN ?", postCommentIDs...),
			qm.And("member_id = ?", blockerId),
		).All(c.Request.Context(), db)

		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 좋아요를 누른 comment_id를 맵으로 저장 (빠른 조회를 위해)
		likedCommentMap := make(map[int64]bool)
		for _, like := range likes {
			likedCommentMap[like.PostCommentID] = true
		}

		// Initialize a slice to hold all comments
		var topLevelComments []PostCommentResponse

		// Add all top-level comments (those without parent comments) to the slice
		for _, comment := range postComments {
			if !comment.IsRecomment.Bool {
				if comment.CreatedAt.Valid {
				}
				// Top-level comment, add to slice
				topLevelComments = append(topLevelComments, PostCommentResponse{
					CommentId:       comment.PostCommentID,
					Content:         comment.Content.String,
					IsRecomment:     comment.IsRecomment.Bool,
					ParentCommentId: comment.ParentPostCommentID.Int64,
					PostId:          comment.PostID,
					MemberId:        comment.MemberID,
					Nickname:        comment.R.Member.Nickname.String,
					CreatedAt:       comment.CreatedAt.Time,
					Likes:           comment.Likes,
					IsLiked:         likedCommentMap[comment.PostCommentID],
					PostRecomments:  []PostCommentResponse{},
				})
			}
		}

		// Add reComments to their respective parent comments in the slice
		for _, comment := range postComments {
			if comment.IsRecomment.Bool {
				// Find the parent comment in the topLevelComments slice and append the recomment
				for i := range topLevelComments {
					if topLevelComments[i].CommentId == comment.ParentPostCommentID.Int64 {
						reComment := PostCommentResponse{
							CommentId:       comment.PostCommentID,
							Content:         comment.Content.String,
							IsRecomment:     comment.IsRecomment.Bool,
							ParentCommentId: comment.ParentPostCommentID.Int64,
							MemberId:        comment.MemberID,
							Nickname:        comment.R.Member.Nickname.String,
							CreatedAt:       comment.CreatedAt.Time,
							PostId:          comment.PostID,
							Likes:           comment.Likes,
							IsLiked:         likedCommentMap[comment.PostCommentID],
							PostRecomments:  []PostCommentResponse{},
						}
						topLevelComments[i].PostRecomments = append(topLevelComments[i].PostRecomments, reComment)
						break
					}
				}
			}
		}

		// Sort reComments by CreatedAt within each top-level comment
		for i := range topLevelComments {
			sort.Slice(topLevelComments[i].PostRecomments, func(j, k int) bool {
				return topLevelComments[i].PostRecomments[j].CreatedAt.Before(topLevelComments[i].PostRecomments[k].CreatedAt)
			})
		}

		// Return comments as part of the response
		pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentResponse{TotalPostCommentCount: int(totalCommentsCount), PostComments: topLevelComments, NextPage: pageInt + 1})
	}
}

type GetPostReCommentResponse struct {
	TotalPostReCommentCount int                   `json:"totalPostReCommentCount"`
	PostReComments          []PostCommentResponse `json:"postReComments"`
	NextPage                int                   `json:"nextPage"`
}

// GetReCommentOnPost 댓글에 대한 대댓글 정보 보기
// @Summary      Retrieve rePostComments for the specified PostCommentId
// @Description  Get rePostComments for a specific comment identified by postCommentId
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postCommentId   path      int  true  "Post Comment ID"
// @Param        page query int false "현재 조회할 게시글 목록의 쪽수. 입력하지 않는다면 기본값인 1쪽을 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetPostReCommentResponse} "Success"
// @Router       /v1/posts/comments/{postCommentId}/recomments [get]
// @Security BearerAuth
func GetReCommentOnPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve commentId from path parameter
		postCommentIdParam := c.Param("postCommentId")
		log.Printf("postCommentIdParam: %s", postCommentIdParam)
		postCommentId, err := strconv.Atoi(postCommentIdParam)
		if err != nil {
			log.Println("Error converting postCommentId:", err) // 변환 실패 시 로그
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postCommentId", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		pageStr := c.DefaultQuery("page", defaultPage)
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil || pageInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		// OFFSET 계산
		offset := (pageInt - 1) * sizeInt

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

		//blocked_member_id 리스트 만들기
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		// Retrieve reComments for the specified commentId
		reComments, err := mysql.PostComments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = post_comment.member_id"),
			qm.Where("post_comment.parent_comment_id = ? and post_comment.deleted_at is null", postCommentId),
			qm.WhereNotIn("post_comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.Limit(sizeInt),
			qm.Offset(offset),
			qm.OrderBy("post_comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(reComments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", GetPostReCommentResponse{TotalPostReCommentCount: 0, PostReComments: []PostCommentResponse{}, NextPage: pageInt})
			return
		}

		// Prepare the final data list directly in the order retrieved
		data := make([]PostCommentResponse, 0, len(reComments))
		for _, recomment := range reComments {
			data = append(data, PostCommentResponse{
				CommentId:       recomment.PostCommentID,
				Content:         recomment.Content.String,
				IsRecomment:     recomment.IsRecomment.Bool,
				ParentCommentId: recomment.ParentPostCommentID.Int64,
				PostId:          recomment.PostID,
				MemberId:        recomment.MemberID,
				Nickname:        recomment.R.Member.Nickname.String,
				Likes:           recomment.Likes,
				CreatedAt:       recomment.CreatedAt.Time,
			})
		}

		// 전체 댓글 수를 가져오는 쿼리
		totalReCommentsCount, err := mysql.PostComments(
			qm.Where("post_comment.parent_comment_id = ? AND post_comment.is_recomment = true AND post_comment.deleted_at IS NULL", postCommentId),
			qm.WhereNotIn("post_comment.member_id NOT IN ?", blockedMemberIds...), // 블랙리스트 제외
		).Count(c.Request.Context(), db)

		postRecommendResponse := GetPostReCommentResponse{
			TotalPostReCommentCount: int(totalReCommentsCount),
			PostReComments:          data,
			NextPage:                pageInt + 1,
		}

		// Return the response with the data list
		pkg.BaseResponse(c, http.StatusOK, "success", postRecommendResponse)
	}
}

type PostCommentReportRequest struct {
	PostCommentId   int64  `json:"postCommentId"`
	Reason          string `json:"reason"`
	SubjectMemberId int64  `json:"subjectMemberId"`
}

type PostCommentReportResponse struct {
	PostCommentReportId   int64  `json:"postReportId"`
	PostCommentId         int64  `json:"postCommentId"`
	Reason                string `json:"reason"`
	SubjectMemberId       int64  `json:"subjectMemberId"`
	PostCommentReporterId int64  `json:"postReporterId"`
}

// ReportPostComment godoc
// @Summary      해당하는 댓글ID를 통해 신고하기
// @Description  해당하는 댓글ID를 통해 신고하기
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        ReportRequest   body      ReportRequest  true  "ReportRequest"
// @Success      200 {object} pkg.BaseResponseStruct{data=PostCommentReportResponse} "성공"
// @Router       /v1/posts/comments/report [post]
// @Security BearerAuth
func ReportPostComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportRequest := &PostCommentReportRequest{}
		if err := c.ShouldBindJSON(&reportRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// Get memberId from context
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		nullReason := null.StringFrom(reportRequest.Reason)

		// 댓글 신고하기
		m := mysql.PostCommentReport{PostCommentID: reportRequest.PostCommentId, ReportReason: nullReason, SubjectMemberID: reportRequest.SubjectMemberId, ReporterMemberID: memberId.(int64)}
		err := m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		reportResponse := PostCommentReportResponse{
			PostCommentReportId:   m.PostCommentReportID,
			PostCommentId:         m.PostCommentID,
			Reason:                m.ReportReason.String,
			SubjectMemberId:       m.SubjectMemberID,
			PostCommentReporterId: m.ReporterMemberID,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", reportResponse)
	}
}

// LikePostComment godoc
// @Summary      해당하는 댓글에 좋아요 누르기
// @Description  해당하는 댓글에 좋아요 누르기
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        commentId   path  int  true  "Comment ID"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v1/posts/comments/{postCommentId}/like [post]
// @Security BearerAuth
func LikePostComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// commentId 가져오기
		postCommentIdParam := c.Param("postCommentId")
		postCommentId, err := strconv.ParseInt(postCommentIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postCommentId", nil)
			return
		}

		// 좋아요 상태 변경 함수
		changeLikeStatus := func(comment *mysql.PostComment, delta int) error {
			comment.Likes += 1
			_, err := comment.Update(c, db, boil.Infer())
			return err
		}

		// 이미 좋아요를 눌렀는지 확인
		postCommentLikes, err := mysql.PostCommentLikes(
			qm.Where("member_id = ? AND post_comment_id = ? AND deleted_at IS NULL", memberId.(int64), postCommentId),
		).One(c.Request.Context(), db)

		// 이미 좋아요를 누른 상태에서 좋아요 취소 요청
		if err == nil {
			postCommentLikes.DeletedAt = null.TimeFrom(time.Now())
			if _, err := postCommentLikes.Update(c.Request.Context(), db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// CommentTable에서 해당 CommentId의 LikeCount를 1 감소시킨다
			postComment, err := mysql.PostComments(
				qm.Where("post_comment_id = ?", postCommentId),
			).One(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			if err := changeLikeStatus(postComment, -1); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", postComment.Likes)
			return
		}

		// 댓글 좋아요 누르기
		like := mysql.PostCommentLike{MemberID: memberId.(int64), PostCommentID: postCommentId}
		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// CommentTable에서 해당 CommentId의 LikeCount를 1 증가시킨다
		postComment, err := mysql.PostComments(
			qm.Where("post_comment_id = ?", postCommentId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if err := changeLikeStatus(postComment, 1); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", postComment.Likes)
		return
	}
}
