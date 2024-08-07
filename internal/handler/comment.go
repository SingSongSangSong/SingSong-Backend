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

// Define the CommentResponse struct
type CommentResponse struct {
	CommentId       int64             `json:"commentId"`
	Content         string            `json:"content"`
	IsRecomment     bool              `json:"isRecomment"`
	ParentCommentId int64             `json:"parentCommentId"`
	SongInfoId      int64             `json:"songId"`
	MemberId        int64             `json:"memberId"`
	Nickname        string            `json:"nickname"`
	CreatedAt       time.Time         `json:"createdAt"`
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

		member, err := mysql.Members(qm.Where("member_id = ?", memberId.(int64))).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 댓글 달기
		nulIsRecomment := null.BoolFrom(commentRequest.IsRecomment)
		nullParentCommentId := null.Int64From(commentRequest.ParentCommentId)
		nullContent := null.StringFrom(commentRequest.Content)
		m := mysql.Comment{MemberID: memberId.(int64), ParentCommentID: nullParentCommentId, SongInfoID: commentRequest.SongInfoId, IsRecomment: nulIsRecomment, Content: nullContent}
		err = m.Insert(c, db, boil.Infer())
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
		// Retrieve memberId from context (for potential blocking features)
		//memberId, exists := c.Get("memberId")
		//if !exists {
		//	pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
		//	return
		//}

		// Retrieve songId from path parameter
		songIdParam := c.Param("songId")
		songId, err := strconv.Atoi(songIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid songId", nil)
			return
		}

		// Retrieve comments for the specified songId
		//comments, err := mysql.Comments(qm.Where("song_info_id = ?", songId), qm.OrderBy("created_at ASC")).All(c, db)
		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.song_info_id = ?", songId),
			qm.OrderBy("comment.created_at ASC"),
		).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Create a map to organize comments and their recomments
		commentMap := make(map[int64]*CommentResponse)

		// First, add all top-level comments (those without parent comments)
		for _, comment := range comments {
			if !comment.IsRecomment.Bool {
				// Top-level comment, add to map
				commentMap[comment.CommentID] = &CommentResponse{
					CommentId:       comment.CommentID,
					Content:         comment.Content.String,
					IsRecomment:     comment.IsRecomment.Bool,
					ParentCommentId: comment.ParentCommentID.Int64,
					SongInfoId:      comment.SongInfoID,
					MemberId:        comment.MemberID,
					Nickname:        comment.R.Member.Nickname.String,
					CreatedAt:       comment.CreatedAt.Time,
					Recomments:      []CommentResponse{},
				}
			}
		}

		// Add recomments to their respective parent comments
		for _, comment := range comments {
			if comment.IsRecomment.Bool {
				// Recomment, add to the parent's Recomments slice
				if parent, exists := commentMap[comment.ParentCommentID.Int64]; exists {
					recomment := CommentResponse{
						CommentId:       comment.CommentID,
						Content:         comment.Content.String,
						IsRecomment:     comment.IsRecomment.Bool,
						ParentCommentId: comment.ParentCommentID.Int64,
						MemberId:        comment.MemberID,
						Nickname:        comment.R.Member.Nickname.String,
						CreatedAt:       comment.CreatedAt.Time,
						SongInfoId:      comment.SongInfoID,
					}
					parent.Recomments = append(parent.Recomments, recomment)
				}
			}
		}

		// Prepare the final data list
		data := make([]CommentResponse, 0, len(commentMap))
		for _, comment := range commentMap {
			data = append(data, *comment)
		}

		// Sort the data slice by CreatedAt timestamp
		sort.Slice(data, func(i, j int) bool {
			return data[i].CreatedAt.Before(data[j].CreatedAt)
		})

		// Sort recomments by CreatedAt within each top-level comment
		for i := range data {
			sort.Slice(data[i].Recomments, func(j, k int) bool {
				return data[i].Recomments[j].CreatedAt.Before(data[i].Recomments[k].CreatedAt)
			})
		}

		// Return comments as part of the response
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
		err := m.Insert(c, db, boil.Infer())
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
// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
// @Router       /comment/like [post]
// @Security BearerAuth
func LikeComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId가져오기

		// 댓글 좋아요 누르기
	}
}
