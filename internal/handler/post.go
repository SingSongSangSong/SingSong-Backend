package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"time"
)

// todo: 이미지 등록

type PostRequest struct {
	Title       string  `json:"title"`
	Content     string  `json:"content"`
	SongInfoIds []int64 `json:"songIds"`
}

type PostIdResponse struct {
	PostId int64 `json:"postId"`
}

// CreatePost godoc
// @Summary      게시글 등록
// @Description  게시글 등록
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        PostRequest   body      PostRequest  true  "PostRequest"
// @Success      200 {object} pkg.BaseResponseStruct{data=PostIdResponse} "성공"
// @Failure      400 "PostRequest가 올바르지 않은 경우, 엑세스 토큰은 유효하지만 사용자 정보가 유효하지 않은 경우, 노래 id가 중복되거나 10개 초과인 경우 400 실패"
// @Failure      401 "엑세스 토큰 검증에 실패했을 경우 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts [post]
// @Security BearerAuth
func CreatePost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postRequest := &PostRequest{}
		if err := c.ShouldBindJSON(postRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// memberId가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		exists, err := mysql.Members(qm.Where("member_id = ?", memberId.(int64))).Exists(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid member", nil)
			return
		}

		if len(postRequest.SongInfoIds) > 10 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - maximum song size is 10", nil)
			return
		}

		if postRequest.SongInfoIds == nil || len(postRequest.SongInfoIds) == 0 {
			post := mysql.Post{
				BoardID:  1, //default board
				MemberID: memberId.(int64),
				Title:    postRequest.Title,
				Content:  null.StringFrom(postRequest.Content),
			}
			err = post.Insert(c.Request.Context(), db, boil.Infer())
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", PostIdResponse{post.PostID})
			return
		}

		// songInfoId 유효성 검사/중복체크 시작
		var songInfoIds []interface{}
		seen := make(map[int64]bool) // 중복 체크를 위한 맵

		for _, songInfoId := range postRequest.SongInfoIds {
			// 이미 처리된 songInfoId가 있으면 에러 반환
			if _, exists := seen[songInfoId]; exists {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - duplicate song info ID detected", nil)
				return
			}
			seen[songInfoId] = true
			songInfoIds = append(songInfoIds, songInfoId)
		}

		count, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).Count(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		if count != int64(len(postRequest.SongInfoIds)) {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - contains invalid song id", nil)
			return
		}

		post := mysql.Post{
			BoardID:  1, //default board
			MemberID: memberId.(int64),
			Title:    postRequest.Title,
			Content:  null.StringFrom(postRequest.Content),
		}
		err = post.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 노래 Batch insert 쿼리 생성
		query := "INSERT INTO post_song (post_id, song_info_id) VALUES "
		var values []interface{}
		for _, songInfoID := range postRequest.SongInfoIds {
			query += "(?, ?),"
			values = append(values, post.PostID, songInfoID)
		}

		// 마지막 콤마 제거
		query = query[:len(query)-1]

		// Batch insert 실행
		_, err = db.ExecContext(c.Request.Context(), query, values...)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", PostIdResponse{post.PostID})
	}
}

type PostDetailsResponse struct {
	PostId      int64        `json:"postId"`
	Title       string       `json:"title"`
	Content     string       `json:"content"`
	Likes       int          `json:"likes"`
	MemberId    int64        `json:"memberId"`
	IsWriter    bool         `json:"isWriter"`
	Nickname    string       `json:"nickname"`
	IsLiked     bool         `json:"isLiked"`
	CreatedAt   string       `json:"createdAt"`
	SongsOnPost []SongOnPost `json:"songs"`
}

type SongOnPost struct {
	SongNumber int    `json:"songNumber"`
	SongName   string `json:"songName"`
	SingerName string `json:"singerName"`
	SongInfoId int64  `json:"songId"`
	Album      string `json:"album"`
	IsMr       bool   `json:"isMr"`
	IsLive     bool   `json:"isLive"`
	MelonLink  string `json:"melonLink"`
}

// GetPost godoc
// @Summary      게시글 하나 상세 조회
// @Description  게시글 하나 상세 조회
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId path string true "postId"
// @Success      200 {object} pkg.BaseResponseStruct{data=PostDetailsResponse} "성공"
// @Failure      400 "postId가 요청에 없는 경우, 해당 게시글이 존재하지 않는 경우 400 실패"
// @Failure      401 "엑세스 토큰 검증에 실패했을 경우 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/{postId} [get]
// @Security BearerAuth
func GetPost(db *sql.DB) gin.HandlerFunc {
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

		one, err := mysql.Posts(
			qm.Load(mysql.PostRels.PostSongs),
			qm.LeftOuterJoin("post_song on post_song.post_id = post.post_id"),
			qm.Where("post.post_id = ? and post.deleted_at is null", postId),
		).One(c.Request.Context(), db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// postId에 해당하는 게시글이 존재하지 않는 경우
				pkg.BaseResponse(c, http.StatusBadRequest, "error - post not found", nil)
				return
			}
			// 기타 데이터베이스 관련 에러
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		isLiked, err := mysql.PostLikes(qm.Where("post_id = ? and member_id = ? and deleted_at is null", postId, memberId)).Exists(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		writer, err := mysql.Members(qm.Where("member_id = ? and deleted_at is null", memberId)).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		var songInfoIds []interface{}
		for _, postSong := range one.R.PostSongs {
			songInfoIds = append(songInfoIds, postSong.SongInfoID)
		}
		all, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		var songsOnPost []SongOnPost
		for _, song := range all {
			songsOnPost = append(songsOnPost, SongOnPost{
				SongNumber: song.SongNumber,
				SongName:   song.SongName,
				SingerName: song.ArtistName,
				SongInfoId: song.SongInfoID,
				Album:      song.Album.String,
				IsMr:       song.IsMR.Bool,
				IsLive:     song.IsLive.Bool,
				MelonLink:  CreateMelonLinkByMelonSongId(song.MelonSongID),
			})
		}

		if len(all) == 0 {
			songsOnPost = []SongOnPost{}
		}

		response := PostDetailsResponse{
			PostId:      one.PostID,
			Title:       one.Title,
			Content:     one.Content.String,
			Likes:       one.Likes,
			MemberId:    one.MemberID,
			IsWriter:    one.MemberID == memberId,
			IsLiked:     isLiked,
			Nickname:    writer.Nickname.String,
			CreatedAt:   one.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			SongsOnPost: songsOnPost,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", response)
	}
}

// DeletePost godoc
// @Summary      게시글 하나 삭제
// @Description  게시글 하나 삭제
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId path string true "postId"
// @Success      200 "성공"
// @Failure      400 "postId가 요청에 없는 경우, 해당 게시글이 존재하지 않는 경우, 게시글 작성자가 아닌 경우 400 실패"
// @Failure      401 "사용자 인증에 실패했을 경우 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/{postId} [delete]
// @Security BearerAuth
func DeletePost(db *sql.DB) gin.HandlerFunc {
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

		one, err := mysql.Posts(
			qm.Where("post.post_id = ? and post.deleted_at is null", postId),
		).One(c.Request.Context(), db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// postId에 해당하는 게시글이 존재하지 않는 경우
				pkg.BaseResponse(c, http.StatusBadRequest, "error - post not found", nil)
				return
			}
			// 기타 데이터베이스 관련 에러
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if one.MemberID != memberId {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - you are not writer", nil)
			return
		}

		_, err = mysql.Posts(qm.Where("post.post_id = ? and post.deleted_at is null", postId)).
			UpdateAll(c.Request.Context(), db, mysql.M{
				"deleted_at": time.Now(),
			})

		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

type postPageResponse struct {
	Posts      []postPreviewResponse `json:"posts"`
	LastCursor int64                 `json:"lastCursor"`
}

type postPreviewResponse struct {
	PostId       int64  `json:"postId"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Nickname     string `json:"nickname"`
	Likes        int    `json:"likes"`
	CommentCount int    `json:"commentCount"`
}

// ListPosts godoc
// @Summary      게시글 전체 조회 (커서 기반 페이징)
// @Description  게시글 전체 조회 (커서 기반 페이징)
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        cursor query int false "마지막에 조회했던 커서의 postId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신 글부터 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=postPageResponse} "성공"
// @Failure      400 "query param 값이 들어왔는데, 숫자가 아니라면 400 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts [get]
func ListPosts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt <= 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		// 페이징 처리된 게시글 가져오기
		posts, err := mysql.Posts(
			qm.Select("DISTINCT post.*"), // 중복된 게시글 제거
			qm.Load(mysql.PostRels.PostComments),
			qm.Load(mysql.PostRels.Member),
			qm.LeftOuterJoin("post_comment on post_comment.post_id = post.post_id"),
			qm.Where("post.post_id < ?", cursorInt),
			qm.OrderBy("post.post_id DESC"),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		previews := make([]postPreviewResponse, 0, sizeInt)

		for _, post := range posts {
			comments := post.R.PostComments
			previews = append(previews, postPreviewResponse{
				PostId:       post.PostID,
				Title:        post.Title,
				Content:      post.Content.String,
				Nickname:     post.R.Member.Nickname.String,
				Likes:        post.Likes,
				CommentCount: len(comments),
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(previews) > 0 {
			lastCursor = previews[len(previews)-1].PostId
		}

		response := postPageResponse{
			Posts:      previews,
			LastCursor: lastCursor,
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type PostReportRequest struct {
	Reason          string `json:"reason"`
	SubjectMemberId int64  `json:"subjectMemberId"`
}

// ReportPost godoc
// @Summary      게시글 신고
// @Description  게시글 신고
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId path string true "postId"
// @Param        PostReportRequest   body      PostReportRequest  true  "PostReportRequest"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Failure      400 "postId param이 잘못 들어왔거나, body 형식이 올바르지 않다면 400 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/{postId}/reports [post]
// @Security BearerAuth
func ReportPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIdStr := c.Param("postId")
		if postIdStr == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find postId in path variable", nil)
			return
		}

		// string을 int64로 변환
		postId, err := strconv.ParseInt(postIdStr, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - postId type invalid", nil)
			return
		}

		reportRequest := &PostReportRequest{}
		if err := c.ShouldBindJSON(&reportRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		nullReason := null.StringFrom(reportRequest.Reason)

		m := mysql.PostReport{PostID: postId, ReportReason: nullReason, SubjectMemberID: reportRequest.SubjectMemberId, ReporterMemberID: memberId.(int64)}
		err = m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

// LikePost godoc
// @Summary      해당하는 게시글에 좋아요 누르기
// @Description  해당하는 게시글에 좋아요 누르기
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId   path  int  true  "Post ID"
// @Success      200 {object} pkg.BaseResponseStruct{data=int} "성공"
// @Router       /v1/posts/{postId}/likes [post]
// @Security BearerAuth
func LikePost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// memberId 가져오기
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		// postId 가져오기
		postIdParam := c.Param("postId")
		postId, err := strconv.ParseInt(postIdParam, 10, 64)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postId", nil)
			return
		}

		// 좋아요 상태 변경 함수
		changeLikeStatus := func(post *mysql.Post, delta int) error {
			post.Likes += delta
			_, err := post.Update(c, db, boil.Infer())
			return err
		}

		// 이미 좋아요를 눌렀는지 확인
		postLikes, err := mysql.PostLikes(
			qm.Where("member_id = ? AND post_id = ? AND deleted_at IS NULL", memberId.(int64), postId),
		).One(c.Request.Context(), db)

		// 이미 좋아요를 누른 상태에서 좋아요 취소 요청
		if err == nil {
			postLikes.DeletedAt = null.TimeFrom(time.Now())
			if _, err := postLikes.Update(c.Request.Context(), db, boil.Infer()); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// Post Table에서 해당 postId의 LikeCount를 1 감소시킨다
			post, err := mysql.Posts(
				qm.Where("post_id = ?", postId),
			).One(c.Request.Context(), db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			if err := changeLikeStatus(post, -1); err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", post.Likes)
			return
		}

		// 게시글 좋아요 누르기
		like := mysql.PostLike{MemberID: memberId.(int64), PostID: postId}
		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// CommentTable에서 해당 CommentId의 LikeCount를 1 증가시킨다
		post, err := mysql.Posts(
			qm.Where("post_id = ?", postId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if err := changeLikeStatus(post, 1); err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", post.Likes)
		return
	}
}

// SearchPosts godoc
// @Summary      게시글 검색 및 조회 (커서 기반 페이징)
// @Description  게시글 검색 및 조회 (커서 기반 페이징)
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword query string true "검색 키워드"
// @Param        cursor query int false "마지막에 조회했던 커서의 postId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 최신 글부터 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=postPageResponse} "성공"
// @Failure      400 "query param 값이 들어왔는데, 비어있다면 400 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/search/posts [get]
func SearchPosts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 검색어를 쿼리 파라미터에서 가져오기
		searchKeyword := c.Query("keyword")
		if searchKeyword == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find keyword in query", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSearchSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		cursorStr := c.DefaultQuery("cursor", "9223372036854775807") //int64 최대값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt <= 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		// 페이징 처리된 게시글 가져오기
		posts, err := mysql.Posts(
			qm.Select("DISTINCT post.*"), // 중복된 게시글 제거
			qm.Load(mysql.PostRels.PostComments),
			qm.Load(mysql.PostRels.Member),
			qm.LeftOuterJoin("post_comment on post_comment.post_id = post.post_id"),
			qm.Where("post.post_id < ? AND (post.title LIKE ? OR post.content LIKE ?)", cursorInt, "%"+searchKeyword+"%", "%"+searchKeyword+"%"),
			qm.OrderBy("post.post_id DESC"),
			qm.Limit(sizeInt),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		previews := make([]postPreviewResponse, 0, sizeInt)

		for _, post := range posts {
			comments := post.R.PostComments
			previews = append(previews, postPreviewResponse{
				PostId:       post.PostID,
				Title:        post.Title,
				Content:      post.Content.String,
				Nickname:     post.R.Member.Nickname.String,
				Likes:        post.Likes,
				CommentCount: len(comments),
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		var lastCursor int64 = 0
		if len(previews) > 0 {
			lastCursor = previews[len(previews)-1].PostId
		}

		response := postPageResponse{
			Posts:      previews,
			LastCursor: lastCursor,
		}

		// 응답 반환
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
