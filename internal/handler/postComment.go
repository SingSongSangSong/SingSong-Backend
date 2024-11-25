package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"errors"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostCommentRequest struct {
	ParentCommentId int64   `json:"parentCommentId"`
	PostId          int64   `json:"postId"`
	Content         string  `json:"content"`
	IsRecomment     bool    `json:"isRecomment"`
	SongInfoIds     []int64 `json:"songIds"`
}

type PostCommentResponse struct {
	PostCommentId       int64        `json:"postCommentId"`
	Content             string       `json:"content"`
	IsRecomment         bool         `json:"isRecomment"`
	ParentPostCommentId int64        `json:"parentPostCommentId"`
	PostId              int64        `json:"postId"`
	MemberId            int64        `json:"memberId"`
	IsWriter            bool         `json:"isWriter"`
	Nickname            string       `json:"nickname"`
	CreatedAt           time.Time    `json:"createdAt"`
	Likes               int          `json:"likes"`
	IsLiked             bool         `json:"isLiked"`
	SongOnPostComment   []SongOnPost `json:"songOnPostComment"`
	PostRecommentCount  int          `json:"postRecommentsCount"`
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
func CommentOnPost(db *sql.DB, firebaseApp *firebase.App) gin.HandlerFunc {
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
			pkg.SendToSentryWithStack(c, fmt.Errorf("memberId not found in context"))
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		member, err := mysql.Members(qm.Where("member_id = ?", memberId.(int64))).One(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// songInfoId 유효성 검사/중복체크 시작
		var songInfoIds []interface{}
		seen := make(map[int64]bool) // 중복 체크를 위한 맵

		// 댓글에 대한 노래 정보 저장
		for _, songInfoId := range commentRequest.SongInfoIds {
			if _, exists := seen[songInfoId]; exists {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - duplicate song info ID detected", nil)
				return
			}
			seen[songInfoId] = true
			songInfoIds = append(songInfoIds, songInfoId)
		}

		// 댓글 달기
		nulIsRecomment := null.BoolFrom(commentRequest.IsRecomment)
		nullParentCommentId := null.Int64From(commentRequest.ParentCommentId)
		nullContent := null.StringFrom(commentRequest.Content)
		postComment := mysql.PostComment{MemberID: memberId.(int64), ParentPostCommentID: nullParentCommentId, PostID: commentRequest.PostId, IsRecomment: nulIsRecomment, Content: nullContent, Likes: 0}
		err = postComment.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(songInfoIds) > 0 {
			count, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).Count(c.Request.Context(), db)
			if err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			if count != int64(len(commentRequest.SongInfoIds)) {
				pkg.BaseResponse(c, http.StatusBadRequest, "error - contains invalid song id", nil)
				return
			}

			// 노래 Batch insert 쿼리 생성
			query := "INSERT INTO post_comment_song (post_comment_id, song_info_id) VALUES "
			var values []interface{}
			for _, songInfoID := range commentRequest.SongInfoIds {
				query += "(?, ?),"
				values = append(values, postComment.PostCommentID, songInfoID)
			}

			// 마지막 콤마 제거
			query = query[:len(query)-1]

			// Batch insert 실행
			_, err = db.ExecContext(c.Request.Context(), query, values...)
			if err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}

		commentResponse := PostCommentResponse{
			PostCommentId:       postComment.PostCommentID,
			ParentPostCommentId: postComment.ParentPostCommentID.Int64,
			PostId:              postComment.PostID,
			Content:             postComment.Content.String,
			IsRecomment:         postComment.IsRecomment.Bool,
			MemberId:            postComment.MemberID,
			IsWriter:            member.MemberID == postComment.MemberID,
			Nickname:            member.Nickname.String,
			CreatedAt:           postComment.CreatedAt.Time,
			Likes:               postComment.Likes,
			SongOnPostComment:   make([]SongOnPost, 0),
			PostRecommentCount:  0,
		}

		if postComment.IsRecomment.Bool { //대댓글인경우
			// 대댓글 달렸다고 알림 보내기
			go NotifyRecommentOnPostComment(db, firebaseApp, memberId.(int64), postComment.ParentPostCommentID.Int64, postComment.PostID, postComment.Content.String)
		} else { //부모댓글인 경우
			// 댓글이 달렸다고 알림 보내기
			go NotifyCommentOnPost(db, firebaseApp, memberId.(int64), commentRequest.PostId, commentRequest.Content)
		}

		// 댓글 달기 성공시 댓글 정보 반환
		pkg.BaseResponse(c, http.StatusOK, "success", commentResponse)
	}
}

type GetPostCommentResponse struct {
	TotalPostCommentCount int                   `json:"totalPostCommentCount"`
	PostComments          []PostCommentResponse `json:"postComments"`
	LastCursor            int64                 `json:"lastCursor"`
}

// PostCommentWithCounts 구조체 정의
type PostCommentWithCounts struct {
	mysql.PostComment             // PostComment 구조체
	NickName          null.String // 닉네임
	ReplyCount        int         // 대댓글 수
	SongInfoIds       null.String // 노래 정보 ID
}

// GetCommentOnPost godoc
// @Summary      Retrieve comments for the specified postId
// @Description  Get comments for a specific post identified by postId with optional page and size query parameters
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId   path     int  true  "Post ID"
// @Param        cursor query int false "마지막에 조회했던 커서의 postCommentId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 먼저 작성된 댓글부터 조회"
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

		cursorStr := c.DefaultQuery("cursor", "0") //int64 최소값
		cursorInt, err := strconv.Atoi(cursorStr)
		if err != nil || cursorInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid cursor parameter", nil)
			return
		}

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.SendToSentryWithStack(c, fmt.Errorf("memberId not found in context"))
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		//차단 유저 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		blockedMemberIds := make([]interface{}, 0, len(blacklists))
		// blockedMemberIds가 빈 슬라이스가 아닐 경우 쿼리에서 사용하기 위해 변환
		blockedMemberIdsPlaceholder := "0" // 기본적으로 차단된 사용자가 없는 경우
		if len(blacklists) > 0 {
			// blockedMemberIds 슬라이스를 콤마로 구분된 문자열로 변환
			ids := make([]string, len(blacklists))
			for i, blacklist := range blacklists {
				ids[i] = fmt.Sprintf("%v", blacklist.BlockedMemberID)
				blockedMemberIds = append(blockedMemberIds, blacklist.BlockedMemberID)
			}
			blockedMemberIdsPlaceholder = strings.Join(ids, ",")
		}

		// 쿼리 작성
		query := fmt.Sprintf(`
			SELECT post_comment.*,
			   member.nickname as nickname,
			   COUNT(replies.post_comment_id) AS reply_count,
			   GROUP_CONCAT(comment_song.song_info_id) AS song_info_ids
			FROM post_comment
			LEFT JOIN member ON member.member_id = post_comment.member_id
			LEFT JOIN post_comment AS replies 
				ON replies.parent_post_comment_id = post_comment.post_comment_id 
				AND replies.is_recomment = true
			    AND replies.deleted_at IS NULL
			LEFT JOIN post_comment_song AS comment_song ON post_comment.post_comment_id = comment_song.post_comment_id
			WHERE post_comment.post_id = ? 
				AND post_comment.deleted_at IS NULL 
				AND post_comment.is_recomment = FALSE
				AND post_comment.post_comment_id > ?
				AND post_comment.member_id NOT IN (%s) -- 블록된 사용자 제외
			GROUP BY post_comment.post_comment_id
			ORDER BY post_comment.created_at ASC
			LIMIT ?
		`, blockedMemberIdsPlaceholder)

		// Query 실행
		rows, err := db.Query(query, postId, cursorInt, sizeInt)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		defer rows.Close()

		// 결과를 담을 구조체 슬라이스 생성
		postComments := make([]PostCommentWithCounts, 0, sizeInt)
		// 조회 결과를 반복하면서 값을 스캔
		for rows.Next() {
			var postComment PostCommentWithCounts
			err := rows.Scan(
				&postComment.PostComment.PostCommentID,
				&postComment.PostComment.PostID,
				&postComment.PostComment.MemberID,
				&postComment.PostComment.Content,
				&postComment.PostComment.Likes,
				&postComment.PostComment.IsRecomment,
				&postComment.PostComment.ParentPostCommentID,
				&postComment.PostComment.CreatedAt,
				&postComment.PostComment.UpdatedAt,
				&postComment.PostComment.DeletedAt,
				&postComment.NickName,
				&postComment.ReplyCount, // 대댓글 수
				&postComment.SongInfoIds,
			)
			if err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			postComments = append(postComments, postComment)
		}

		// 전체 댓글 수를 가져오는 쿼리
		totalCommentsCount, err := mysql.PostComments(
			qm.Where("post_comment.post_id = ? AND post_comment.deleted_at IS NULL", postId),
			qm.WhereNotIn("post_comment.member_id NOT IN ?", blockedMemberIds...), // 블랙리스트 제외
		).Count(c.Request.Context(), db)

		// 결과가 없는 경우 빈 리스트 반환
		if len(postComments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentResponse{TotalPostCommentCount: int(totalCommentsCount), PostComments: []PostCommentResponse{}, LastCursor: int64(cursorInt)}) //마지막(가장 최근 id) 커서값, 없으면 0
			return
		}

		var songInfoIds []interface{}
		// comment_id들만 추출
		postCommentIDs := make([]interface{}, len(postComments))
		// 중복 제거를 위한 맵 생성
		uniqueIds := make(map[int]bool)
		for i, postComment := range postComments {
			postCommentIDs[i] = postComment.PostCommentID
			if postComment.SongInfoIds.Valid {
				rawData := postComment.SongInfoIds.String
				// 앞뒤 대괄호 제거
				rawData = strings.Trim(rawData, "[]")
				// 쉼표 기준으로 split
				idStrings := strings.Split(rawData, ",")

				for _, idStr := range idStrings {
					// 문자열을 int로 변환
					id, err := strconv.Atoi(strings.TrimSpace(idStr))
					if err == nil {
						if !uniqueIds[id] { // 중복 확인
							uniqueIds[id] = true
							songInfoIds = append(songInfoIds, id) // []interface{}에 추가
						}
					}
				}
			}
		}

		all, err := mysql.SongInfos(qm.WhereIn("song_info_id IN ?", songInfoIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// songInfoId와 songInfo를 맵으로 저장 (빠른 조회를 위해)
		songInfoMap := make(map[int64]mysql.SongInfo)
		for _, songInfo := range all {
			songInfoMap[songInfo.SongInfoID] = *songInfo
		}

		// Initialize a map to hold SongOnPost for each postCommentID
		songOnPostMap := make(map[int64][]SongOnPost)

		// Iterate over each postComment
		for _, postComment := range postComments {
			if postComment.SongInfoIds.Valid {
				rawData := postComment.SongInfoIds.String
				// Remove brackets and split by comma
				rawData = strings.Trim(rawData, "[]")
				idStrings := strings.Split(rawData, ",")

				// Create a list of SongOnPost for this postComment
				songsOnPost := make([]SongOnPost, 0, len(idStrings))
				for _, idStr := range idStrings {
					// Convert string to int
					id, err := strconv.Atoi(strings.TrimSpace(idStr))
					if err == nil {
						// Check if songInfo exists in the map
						if songInfo, exists := songInfoMap[int64(id)]; exists {
							// Create SongOnPost object
							songOnPost := SongOnPost{
								SongNumber:        songInfo.SongNumber,
								SongName:          songInfo.SongName,
								SingerName:        songInfo.ArtistName,
								SongInfoId:        songInfo.SongInfoID,
								Album:             songInfo.Album.String,
								IsMr:              songInfo.IsMR.Bool,
								IsLive:            songInfo.IsLive.Bool,                               // Set according to your logic
								MelonLink:         CreateMelonLinkByMelonSongId(songInfo.MelonSongID), // Set according to your logic
								LyricsYoutubeLink: songInfo.LyricsVideoLink.String,
								TJYoutubeLink:     songInfo.TJYoutubeLink.String,
								LyricsVideoID:     ExtractVideoID(songInfo.LyricsVideoLink.String),
								TJVideoID:         ExtractVideoID(songInfo.TJYoutubeLink.String),
							}
							// Add to the list of songs for this postComment
							songsOnPost = append(songsOnPost, songOnPost)
						}

					}
				}
				// Map postCommentID to the list of songs
				songOnPostMap[postComment.PostCommentID] = songsOnPost
			}
		}

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.PostCommentLikes(
			qm.WhereIn("post_comment_id IN ?", postCommentIDs...),
			qm.And("member_id = ?", blockerId),
		).All(c.Request.Context(), db)

		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 좋아요를 누른 comment_id를 맵으로 저장 (빠른 조회를 위해)
		likedCommentMap := make(map[int64]bool)
		for _, like := range likes {
			likedCommentMap[like.PostCommentID] = true
		}

		// Initialize a slice to hold all comments
		topLevelComments := make([]PostCommentResponse, 0, sizeInt)

		// Add all top-level comments (those without parent comments) to the slice
		for _, comment := range postComments {
			if !comment.IsRecomment.Bool {
				if songOnPostMap[comment.PostCommentID] == nil {
					songOnPostMap[comment.PostCommentID] = []SongOnPost{}
				}
				// Top-level comment, add to slice
				topLevelComments = append(topLevelComments, PostCommentResponse{
					PostCommentId:       comment.PostCommentID,
					Content:             comment.Content.String,
					IsRecomment:         comment.IsRecomment.Bool,
					ParentPostCommentId: comment.ParentPostCommentID.Int64,
					PostId:              comment.PostID,
					MemberId:            comment.MemberID,
					IsWriter:            comment.MemberID == blockerId,
					Nickname:            comment.NickName.String,
					CreatedAt:           comment.CreatedAt.Time,
					Likes:               comment.Likes,
					IsLiked:             likedCommentMap[comment.PostCommentID],
					SongOnPostComment:   songOnPostMap[comment.PostCommentID], // Add the song list here
					PostRecommentCount:  comment.ReplyCount,
				})
			}
		}

		// 다음 페이지를 위한 커서 값 설정
		lastCursor := int64(cursorInt)
		if len(topLevelComments) > 0 {
			lastCursor = topLevelComments[len(topLevelComments)-1].PostCommentId
		}

		// Return comments as part of the response
		pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentResponse{TotalPostCommentCount: int(totalCommentsCount), PostComments: topLevelComments, LastCursor: lastCursor})
	}
}

type GetPostReCommentResponse struct {
	PostReComments []PostCommentResponse `json:"postReComments"`
	LastCursor     int64                 `json:"lastCursor"`
}

// GetReCommentOnPost 댓글에 대한 대댓글 정보 보기
// @Summary      Retrieve rePostComments for the specified PostCommentId
// @Description  Get rePostComments for a specific comment identified by postCommentId
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postCommentId path string true "postCommentId"
// @Param        cursor query int false "마지막에 조회했던 커서의 postCommentId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 먼저 작성된 댓글부터 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetPostReCommentResponse} "Success"
// @Router       /v1/posts/comments/{postCommentId}/recomments [get]
// @Security BearerAuth
func GetReCommentOnPost(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve commentId from path parameter
		postCommentIdParam := c.Param("postCommentId")
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

		cursorStr := c.DefaultQuery("cursor", "0") //int64 최소값
		cursorInt, err := strconv.Atoi(cursorStr)
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
			pkg.SendToSentryWithStack(c, err)
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
			qm.Where("post_comment.parent_post_comment_id = ? and post_comment.deleted_at is null", postCommentId),
			qm.Where("post_comment.post_comment_id > ?", cursorInt),
			qm.WhereNotIn("post_comment.member_id not IN ?", blockedMemberIds...), // 블랙리스트 제외
			qm.Limit(sizeInt),
			qm.OrderBy("post_comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(reComments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", GetPostReCommentResponse{PostReComments: []PostCommentResponse{}, LastCursor: int64(cursorInt)})
			return
		}

		// comment_id들만 추출
		postCommentIDs := make([]interface{}, len(reComments))
		for i, postComment := range reComments {
			postCommentIDs[i] = postComment.PostCommentID
		}

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.PostCommentLikes(
			qm.WhereIn("post_comment_id IN ?", postCommentIDs...),
			qm.And("member_id = ?", blockerId),
		).All(c.Request.Context(), db)

		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 좋아요를 누른 comment_id를 맵으로 저장 (빠른 조회를 위해)
		likedCommentMap := make(map[int64]bool)
		for _, like := range likes {
			likedCommentMap[like.PostCommentID] = true
		}

		// Prepare the final data list directly in the order retrieved
		data := make([]PostCommentResponse, 0, len(reComments))
		for _, recomment := range reComments {
			data = append(data, PostCommentResponse{
				PostCommentId:       recomment.PostCommentID,
				Content:             recomment.Content.String,
				IsRecomment:         recomment.IsRecomment.Bool,
				ParentPostCommentId: recomment.ParentPostCommentID.Int64,
				PostId:              recomment.PostID,
				MemberId:            recomment.MemberID,
				IsWriter:            recomment.MemberID == blockerId,
				Nickname:            recomment.R.Member.Nickname.String,
				Likes:               recomment.Likes,
				IsLiked:             likedCommentMap[recomment.PostCommentID],
				CreatedAt:           recomment.CreatedAt.Time,
			})
		}

		// 다음 페이지를 위한 커서 값 설정
		lastCursor := int64(cursorInt)
		if len(data) > 0 {
			lastCursor = data[len(data)-1].PostCommentId
		}

		postRecommendResponse := GetPostReCommentResponse{
			PostReComments: data,
			LastCursor:     lastCursor,
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
			pkg.SendToSentryWithStack(c, err)
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
// @Param        postCommentId path string true "postCommentId"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v1/posts/comments/{postCommentId}/like [post]
// @Security BearerAuth
func LikePostComment(db *sql.DB, firebaseApp *firebase.App) gin.HandlerFunc {
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
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			// CommentTable에서 해당 CommentId의 LikeCount를 1 감소시킨다
			postComment, err := mysql.PostComments(
				qm.Where("post_comment_id = ?", postCommentId),
			).One(c.Request.Context(), db)
			if err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			if err := changeLikeStatus(postComment, -1); err != nil {
				pkg.SendToSentryWithStack(c, err)
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}

			pkg.BaseResponse(c, http.StatusOK, "success", postComment.Likes)
			return
		}

		// 댓글 좋아요 누르기
		like := mysql.PostCommentLike{MemberID: memberId.(int64), PostCommentID: postCommentId}
		if err := like.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// CommentTable에서 해당 CommentId의 LikeCount를 1 증가시킨다
		postComment, err := mysql.PostComments(
			qm.Where("post_comment_id = ?", postCommentId),
		).One(c.Request.Context(), db)
		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if err := changeLikeStatus(postComment, 1); err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		go NotifyLikeOnPostComment(db, firebaseApp, memberId.(int64), postCommentId, postComment.PostID, postComment.Content.String)

		pkg.BaseResponse(c, http.StatusOK, "success", postComment.Likes)
		return
	}
}

// DeletePostComment godoc
// @Summary      게시글 댓글 하나 삭제
// @Description  게시글 댓글 하나 삭제
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postCommentId path string true "postCommentId"
// @Success      200 "성공"
// @Failure      400 "postCommentId 요청에 없는 경우, 해당 댓글이 존재하지 않는 경우, 댓글 작성자가 아닌 경우 400 실패"
// @Failure      401 "사용자 인증에 실패했을 경우 401 실패"
// @Failure      500 "서버 에러일 경우 500 실패"
// @Router       /v1/posts/comments/{postCommentId} [delete]
// @Security BearerAuth
func DeletePostComment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postCommentId := c.Param("postCommentId")
		if postCommentId == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find postCommentId in path variable", nil)
			return
		}

		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.SendToSentryWithStack(c, fmt.Errorf("memberId not found in context"))
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		one, err := mysql.PostComments(
			qm.Where("post_comment.post_comment_id = ? and post_comment.deleted_at is null", postCommentId),
		).One(c.Request.Context(), db)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				// postId에 해당하는 게시글이 존재하지 않는 경우
				pkg.BaseResponse(c, http.StatusBadRequest, "error - postComment not found", nil)
				return
			}
			// 기타 데이터베이스 관련 에러
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if one.MemberID != memberId {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - you are not writer", nil)
			return
		}

		_, err = mysql.PostComments(qm.Where("post_comment.post_comment_id = ? and post_comment.deleted_at is null", postCommentId)).
			UpdateAll(c.Request.Context(), db, mysql.M{
				"deleted_at": time.Now(),
			})

		if err != nil {
			pkg.SendToSentryWithStack(c, err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}
