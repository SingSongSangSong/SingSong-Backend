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

		if len(postRequest.SongInfoIds) == 0 {
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
	PostId         int64        `json:"postId"`
	Title          string       `json:"title"`
	Content        string       `json:"content"`
	Likes          int          `json:"likes"`
	IsWriter       bool         `json:"isWriter"`
	WriterNickname string       `json:"writerNickname"`
	IsLiked        bool         `json:"isLiked"`
	CreatedAt      string       `json:"createdAt"`
	SongsOnPost    []SongOnPost `json:"songs"`
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
			PostId:         one.PostID,
			Title:          one.Title,
			Content:        one.Content.String,
			Likes:          one.Likes,
			IsWriter:       one.MemberID == memberId,
			IsLiked:        isLiked,
			WriterNickname: writer.Nickname.String,
			CreatedAt:      one.CreatedAt.Time.Format("2006-01-02 15:04:05"),
			SongsOnPost:    songsOnPost,
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
