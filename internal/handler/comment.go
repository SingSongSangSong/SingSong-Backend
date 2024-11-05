package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	firebase "firebase.google.com/go/v4"
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

// CommentResponse Define the CommentResponse struct
type CommentResponse struct {
	CommentId       int64             `json:"commentId"`
	Content         string            `json:"content"`
	IsRecomment     bool              `json:"isRecomment"`
	ParentCommentId int64             `json:"parentCommentId"`
	SongInfoId      int64             `json:"songId"`
	MemberId        int64             `json:"memberId"`
	IsWriter        bool              `json:"isWriter"`
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
// @Router       /v1/comment [post]
// @Security BearerAuth
func CommentOnSong(db *sql.DB, firebaseApp *firebase.App) gin.HandlerFunc {
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
			IsWriter:        member.MemberID == m.MemberID,
			Nickname:        member.Nickname.String,
			CreatedAt:       m.CreatedAt.Time,
			Likes:           m.Likes.Int,
			Recomments:      []CommentResponse{},
		}

		if m.IsRecomment.Bool { //대댓글인경우
			// 대댓글 달렸다고 알림 보내기
			go NotifyRecommentOnSongComment(db, firebaseApp, m.ParentCommentID.Int64, m.SongInfoID, m.Content.String)
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
// @Router       /v1/comment/{songId} [get]
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
		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.song_info_id = ? and comment.deleted_at is null", songId),
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
					IsWriter:        comment.MemberID == blockerId,
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
							IsWriter:        comment.MemberID == blockerId,
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
// @Router       /v1/comment/recomment/{commentId} [get]
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
			qm.Where("comment.parent_comment_id = ? and comment.deleted_at is null", commentId),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.OrderBy("comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

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
				IsWriter:        recomment.MemberID == blockerId,
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
// @Router       /v1/comment/report [post]
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
// @Router       /v1/comment/{commentId}/like [post]
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

type LatestCommentResponse struct {
	CommentId       int64               `json:"commentId"`
	Content         string              `json:"content"`
	IsRecomment     bool                `json:"isRecomment"`
	ParentCommentId int64               `json:"parentCommentId"`
	MemberId        int64               `json:"memberId"`
	Nickname        string              `json:"nickname"`
	CreatedAt       time.Time           `json:"createdAt"`
	Likes           int                 `json:"likes"`
	IsLiked         bool                `json:"isLiked"`
	Song            SongOfLatestComment `json:"song"`
}

type SongOfLatestComment struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
	IsLive     bool   `json:"isLive"`
	MelonLink  string `json:"melonLink"`
}

// GetLatestComments 홈화면 최신 댓글 가져오기
// @Summary      홈화면 최신 댓글 가져오기
// @Description  홈화면 최신 댓글 가져오기. 쿼리 파라미터인 size를 별도로 지정하지 않으면 default size = 5
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        size   query      int  false  "size"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]LatestCommentResponse} "Success"
// @Router       /v1/comment/latest [get]
// @Security BearerAuth
func GetLatestComments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		sizeValue := c.Query("size")
		if sizeValue == "" {
			sizeValue = "5" //default value
		}

		size, err := strconv.Atoi(sizeValue)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot convert size to int", nil)
			return
		}

		//블랙리스트 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", memberId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.deleted_at is null"),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.OrderBy("comment_id DESC"),                                    // created_at 기준으로 최신 순 정렬
			qm.Limit(size),                                                   // 최신 size개의 댓글만 가져옴
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

		response := make([]LatestCommentResponse, 0, size)

		for _, comment := range comments {
			song := songMap[comment.SongInfoID]
			response = append(response, LatestCommentResponse{
				CommentId:       comment.CommentID,
				Content:         comment.Content.String,
				IsRecomment:     comment.IsRecomment.Bool,
				ParentCommentId: comment.ParentCommentID.Int64,
				MemberId:        comment.MemberID,
				Nickname:        comment.R.Member.Nickname.String,
				CreatedAt:       comment.CreatedAt.Time,
				Likes:           comment.Likes.Int,
				IsLiked:         likedCommentMap[comment.CommentID],
				Song: SongOfLatestComment{
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

		// 성공 응답
		pkg.BaseResponse(c, http.StatusOK, "success", response)
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
// @Router       /v1/comment/my [get]
// @Security BearerAuth
func GetMySongComment(db *sql.DB) gin.HandlerFunc {
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
		cursorInt, err := strconv.ParseInt(cursorStr, 10, 64)
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

// DeleteComment godoc
// @Summary      해당하는 댓글 삭제하기
// @Description  해당하는 댓글 삭제하기
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        commentId   path  int  true  "commentId"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v1/comment/{commentId} [delete]
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

type CommentWithRecommentsCountResponse struct {
	CommentId       int64     `json:"commentId"`
	Content         string    `json:"content"`
	IsRecomment     bool      `json:"isRecomment"`
	ParentCommentId int64     `json:"parentCommentId"`
	SongInfoId      int64     `json:"songId"`
	MemberId        int64     `json:"memberId"`
	Nickname        string    `json:"nickname"`
	CreatedAt       time.Time `json:"createdAt"`
	Likes           int       `json:"likes"`
	IsLiked         bool      `json:"isLiked"`
	RecommentsCount int       `json:"recommentsCount"`
}

// GetHotComment godoc
// @Summary      특정 노래의 핫 댓글 가져오기(현재 핫 댓글 조건: 따봉이 5개이상 박혀있는 것중에 따봉 가장 높은거)
// @Description  특정 노래의 핫 댓글 가져오기. 댓글이 없으면 data가 null로 갑니다. 기본값은 댓글 1개인데, size 쿼리 조절해서 더 가져올수 있어요.
// @Tags         Comment
// @Accept       json
// @Produce      json
// @Param        size query string false "조회할 hot 댓글의 개수. 입력하지 않는다면 기본값은 1"
// @Param        songId path string true "songId"
// @Success      200 {object} pkg.BaseResponseStruct{data=[]CommentWithRecommentsCountResponse} "성공"
// @Router       /v1/songs/{songId}/comments/hot [get]
// @Security BearerAuth
func GetHotCommentOfSong(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		songInfoId := c.Param("songId")
		if songInfoId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songId in path variable", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", "1")
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		//차단 유저 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", memberId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		//blocked_member_id 리스트 만들기
		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		for _, blacklist := range blacklists {
			blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
		}

		comments, err := mysql.Comments(
			qm.Load(mysql.CommentRels.Member),
			qm.LeftOuterJoin("member on member.member_id = comment.member_id"),
			qm.Where("comment.song_info_id = ?", songInfoId),
			qm.Where("comment.deleted_at is null"),
			qm.Where("comment.likes > 4"), // todo: 핫 댓글 조건
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...),
			qm.OrderBy("likes desc"),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 댓글이 없다면 빈 리스트 반환
		if len(comments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "ok", []CommentWithRecommentsCountResponse{})
			return
		}

		// 댓글 ID 리스트 생성
		commentIDs := make([]interface{}, 0, len(comments))
		for _, comment := range comments {
			commentIDs = append(commentIDs, comment.CommentID)
		}

		// 모든 댓글의 좋아요 여부를 한 번에 조회
		likesMap := make(map[int64]bool)
		likedComments, err := mysql.CommentLikes(
			qm.WhereIn("comment_id IN ?", commentIDs...),
			qm.And("member_id = ?", memberId),
			qm.Where("deleted_at is null"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		for _, likedComment := range likedComments {
			likesMap[likedComment.CommentID] = true
		}

		// 모든 댓글의 RecommentsCount를 한 번에 조회
		recomments, err := mysql.Comments(
			qm.WhereIn("parent_comment_id IN ?", commentIDs...),
			qm.WhereNotIn("comment.member_id not IN ?", blockedMemberIds...),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		recommentsCountMap := make(map[int64]int)
		for _, recomment := range recomments {
			if recomment.ParentCommentID.Valid {
				recommentsCountMap[recomment.ParentCommentID.Int64]++
			}
		}

		// 댓글 리스트 생성
		response := make([]CommentWithRecommentsCountResponse, 0, sizeInt)
		for _, comment := range comments {
			response = append(response, CommentWithRecommentsCountResponse{
				CommentId:       comment.CommentID,
				Content:         comment.Content.String,
				IsRecomment:     comment.IsRecomment.Bool,
				ParentCommentId: comment.ParentCommentID.Int64,
				SongInfoId:      comment.SongInfoID,
				MemberId:        comment.MemberID,
				Nickname:        comment.R.Member.Nickname.String,
				CreatedAt:       comment.CreatedAt.Time,
				Likes:           comment.Likes.Int,
				IsLiked:         likesMap[comment.CommentID],
				RecommentsCount: recommentsCountMap[comment.CommentID],
			})
		}

		// 댓글 리스트 응답
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
