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

type CommentRequest struct {
	ParentCommentId int64  `json:"parentCommentId"`
	SongInfoId      int64  `json:"songId"`
	Content         string `json:"content"`
	IsRecomment     bool   `json:"isRecomment"`
}

// CommentResponse Define the CommentResponse struct
type CommentResponse struct {
	CommentId       int64             `json:"commentId"`
	Content         string            `json:"content"`
	IsRecomment     bool              `json:"isRecomment"`
	ParentCommentId int64             `json:"parentCommentId"`
	SongInfoId      int64             `json:"songId"`
	MemberId        int64             `json:"memberId"`
	Nickname        string            `json:"nickname"`
	CreatedAt       time.Time         `json:"createdAt"`
	Likes           int               `json:"likes"`
	IsLiked         bool              `json:"isLiked"`
	Recomments      []CommentResponse `json:"recomments"`
}

// CommentOnSong godoc
// @Summary      SongId에 댓글 달기
// @Description  SongId에 댓글 달기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        CommentRequest   body      CommentRequest  true  "commentRequest"
// @Success      200 {object} pkg.BaseResponseStruct{data=CommentResponse} "성공"
// @Router       /comment [post]
// @Security BearerAuth
func CommentOnSong(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CommentRequest 받기
		commentRequest := &CommentRequest{}
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
		m := mysql.Comment{MemberID: memberId.(int64), ParentCommentID: nullParentCommentId, SongInfoID: commentRequest.SongInfoId, IsRecomment: nulIsRecomment, Content: nullContent, Likes: null.IntFrom(0)}
		err = m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		commentResponse := CommentResponse{
			CommentId:       m.CommentID,
			ParentCommentId: m.ParentCommentID.Int64,
			SongInfoId:      m.SongInfoID,
			Content:         m.Content.String,
			IsRecomment:     m.IsRecomment.Bool,
			MemberId:        m.MemberID,
			Nickname:        member.Nickname.String,
			CreatedAt:       m.CreatedAt.Time,
			Likes:           m.Likes.Int,
			Recomments:      []CommentResponse{},
		}

		// 댓글 달기 성공시 댓글 정보 반환
		pkg.BaseResponse(c, http.StatusOK, "success", commentResponse)
	}
}

// GetCommentOnSong godoc
// @Summary      Retrieve comments for the specified SongId
// @Description  Get comments for a specific song identified by songId
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        songId   path      int  true  "Song ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]CommentResponse} "Success"
// @Router       /comment/{songId} [get]
// @Security BearerAuth
func GetCommentOnSong(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve songId from path parameter
		songIdParam := c.Param("songId")
		songId, err := strconv.Atoi(songIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid songId", nil)
			return
		}

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}
		log.Println("blockerId: ", blockerId)

		//차단 유저 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		log.Printf("blacklists: %d", len(blacklists))

		//blocked_member_id 리스트 만들기
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		// Retrieve comments for the specified songId
		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.song_info_id = ?", songId),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.OrderBy("comment.created_at DESC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(comments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", []CommentResponse{})
			return
		}

		// comment_id들만 추출
		commentIDs := make([]interface{}, len(comments))
		for i, comment := range comments {
			commentIDs[i] = comment.CommentID
		}

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.CommentLikes(
			qm.WhereIn("comment_id IN ?", commentIDs...),
			qm.And("member_id = ?", blockerId),
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

		// Initialize a slice to hold all comments
		var topLevelComments []CommentResponse

		// Add all top-level comments (those without parent comments) to the slice
		for _, comment := range comments {
			if !comment.IsRecomment.Bool {
				// Top-level comment, add to slice
				topLevelComments = append(topLevelComments, CommentResponse{
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
					Recomments:      []CommentResponse{},
				})
			}
		}

		// Add reComments to their respective parent comments in the slice
		for _, comment := range comments {
			if comment.IsRecomment.Bool {
				// Find the parent comment in the topLevelComments slice and append the recomment
				for i := range topLevelComments {
					if topLevelComments[i].CommentId == comment.ParentCommentID.Int64 {
						reComment := CommentResponse{
							CommentId:       comment.CommentID,
							Content:         comment.Content.String,
							IsRecomment:     comment.IsRecomment.Bool,
							ParentCommentId: comment.ParentCommentID.Int64,
							MemberId:        comment.MemberID,
							Nickname:        comment.R.Member.Nickname.String,
							CreatedAt:       comment.CreatedAt.Time,
							SongInfoId:      comment.SongInfoID,
							Likes:           comment.Likes.Int,
							IsLiked:         likedCommentMap[comment.CommentID],
							Recomments:      []CommentResponse{},
						}
						topLevelComments[i].Recomments = append(topLevelComments[i].Recomments, reComment)
						break
					}
				}
			}
		}

		// Sort reComments by CreatedAt within each top-level comment
		for i := range topLevelComments {
			sort.Slice(topLevelComments[i].Recomments, func(j, k int) bool {
				return topLevelComments[i].Recomments[j].CreatedAt.Before(topLevelComments[i].Recomments[k].CreatedAt)
			})
		}

		// Return comments as part of the response
		pkg.BaseResponse(c, http.StatusOK, "success", topLevelComments)
	}
}

// GetReCommentOnSong 댓글에 대한 대댓글 정보 보기
// @Summary      Retrieve reComments for the specified CommentId
// @Description  Get reComments for a specific comment identified by commentId
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        commentId   path      int  true  "Comment ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]CommentResponse} "Success"
// @Router       /comment/recomment/{commentId} [get]
// @Security BearerAuth
func GetReCommentOnSong(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve commentId from path parameter
		commentIdParam := c.Param("commentId")
		commentId, err := strconv.Atoi(commentIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}
		log.Printf("commentId: %d", commentId)

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
		reComments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.parent_comment_id = ?", commentId),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.OrderBy("comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		log.Printf("recomments: %d", len(reComments))

		// Prepare the final data list directly in the order retrieved
		data := make([]CommentResponse, 0, len(reComments))
		for _, recomment := range reComments {
			data = append(data, CommentResponse{
				CommentId:       recomment.CommentID,
				Content:         recomment.Content.String,
				IsRecomment:     recomment.IsRecomment.Bool,
				ParentCommentId: recomment.ParentCommentID.Int64,
				SongInfoId:      recomment.SongInfoID,
				MemberId:        recomment.MemberID,
				Nickname:        recomment.R.Member.Nickname.String,
				Likes:           recomment.Likes.Int,
				CreatedAt:       recomment.CreatedAt.Time,
			})
		}

		// Return the response with the data list
		pkg.BaseResponse(c, http.StatusOK, "success", data)
	}
}

type ReportRequest struct {
	CommentId       int64  `json:"commentId"`
	Reason          string `json:"reason"`
	SubjectMemberId int64  `json:"subjectMemberId"`
}

type ReportResponse struct {
	ReportId        int64  `json:"reportId"`
	CommentId       int64  `json:"commentId"`
	Reason          string `json:"reason"`
	SubjectMemberId int64  `json:"subjectMemberId"`
	ReporterId      int64  `json:"reporterId"`
}

// ReportComment godoc
// @Summary      해당하는 댓글ID를 통해 신고하기
// @Description  해당하는 댓글ID를 통해 신고하기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        ReportRequest   body      ReportRequest  true  "ReportRequest"
// @Success      200 {object} pkg.BaseResponseStruct{data=ReportResponse} "성공"
// @Router       /comment/report [post]
// @Security BearerAuth
func ReportComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		reportRequest := &ReportRequest{}
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
		m := mysql.Report{CommentID: reportRequest.CommentId, ReportReason: nullReason, SubjectMemberID: reportRequest.SubjectMemberId, ReporterMemberID: memberId.(int64)}
		err := m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		reportResponse := ReportResponse{
			ReportId:        m.ReportID,
			CommentId:       m.CommentID,
			Reason:          m.ReportReason.String,
			SubjectMemberId: m.SubjectMemberID,
			ReporterId:      m.ReporterMemberID,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", reportResponse)
	}
}

// LikeComment godoc
// @Summary      해당하는 댓글에 좋아요 누르기
// @Description  해당하는 댓글에 좋아요 누르기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        commentId   path  int  true  "Comment ID"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /comment/{commentId}/like [post]
// @Security BearerAuth
func LikeComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// commentId 가져오기
		commentIdParam := c.Param("commentId")
		commentId, err := strconv.ParseInt(commentIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid commentId", nil)
			return
		}

		// 좋아요 상태 변경 함수
		changeLikeStatus := func(comment *mysql.Comment, delta int) error {
			if comment.Likes.Valid {
				comment.Likes.Int += delta
			} else {
				comment.Likes = null.IntFrom(delta)
			}
			_, err := comment.Update(c, db, boil.Infer())
			return err
		}

		// 이미 좋아요를 눌렀는지 확인
		commentLikes, err := mysql.CommentLikes(
			qm.Where("member_id = ? AND comment_id = ? AND deleted_at IS NULL", memberId.(int64), commentId),
		).One(c.Request.Context(), db)

		// 이미 좋아요를 누른 상태에서 좋아요 취소 요청
		if err == nil {
			commentLikes.DeletedAt = null.TimeFrom(time.Now())
			if _, err := commentLikes.Update(c.Request.Context(), db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// CommentTable에서 해당 CommentId의 LikeCount를 1 감소시킨다
			comment, err := mysql.Comments(
				qm.Where("comment_id = ?", commentId),
			).One(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			if err := changeLikeStatus(comment, -1); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", comment.Likes.Int)
			return
		}

		// 댓글 좋아요 누르기
		like := mysql.CommentLike{MemberID: memberId.(int64), CommentID: commentId}
		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// CommentTable에서 해당 CommentId의 LikeCount를 1 증가시킨다
		comment, err := mysql.Comments(
			qm.Where("comment_id = ?", commentId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if err := changeLikeStatus(comment, 1); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", comment.Likes.Int)
		return
	}
}
