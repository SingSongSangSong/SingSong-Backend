package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PostCommentV2Response struct {
	PostCommentId       int64                   `json:"postCommentId"`
	Content             string                  `json:"content"`
	IsRecomment         bool                    `json:"isRecomment"`
	ParentPostCommentId int64                   `json:"parentPostCommentId"`
	PostId              int64                   `json:"postId"`
	MemberId            int64                   `json:"memberId"`
	Nickname            string                  `json:"nickname"`
	CreatedAt           time.Time               `json:"createdAt"`
	Likes               int                     `json:"likes"`
	IsLiked             bool                    `json:"isLiked"`
	SongOnPostComment   []SongOnPost            `json:"songOnPostComment"`
	PostRecommentCount  int                     `json:"postRecommentsCount"`
	PostRecomments      []PostCommentV2Response `json:"postRecomments"`
}

type GetPostCommentV2Response struct {
	TotalPostCommentCount int                     `json:"totalPostCommentCount"`
	PostComments          []PostCommentV2Response `json:"postComments"`
	LastCursor            int64                   `json:"lastCursor"`
}

// PostCommentV2WithCounts 구조체 정의
type PostCommentV2WithCounts struct {
	mysql.PostComment             // PostComment 구조체
	NickName          null.String // 닉네임
	ReplyCount        int         // 대댓글 수
	SongInfoIds       null.String // 노래 정보 ID
}

// GetCommentOnPostV2 godoc
// @Summary      Retrieve comments for the specified postId (Version. 2)
// @Description  Get comments for a specific post identified by postId with optional page and size query parameters (Version 2)
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        postId   path     int  true  "Post ID"
// @Param        cursor query int false "마지막에 조회했던 커서의 postCommentId(이전 요청에서 lastCursor값을 주면 됨), 없다면 default로 가장 먼저 작성된 댓글부터 조회"
// @Param        size query int false "한번에 조회할 게시글 개수. 입력하지 않는다면 기본값인 20개씩 조회"
// @Success      200 {object} pkg.BaseResponseStruct{data=GetPostCommentV2Response} "Success"
// @Router       /v2/posts/{postId}/comments [get]
// @Security BearerAuth
func GetCommentOnPostV2(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve songId from path parameter
		postIdParam := c.Param("postId")
		postId, err := strconv.Atoi(postIdParam)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid postCommentId", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", "10")
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

		//차단 유저 제외
		blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
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
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			postComments = append(postComments, postComment)
		}

		postCommentIds := make([]interface{}, len(postComments))
		for i, postComment := range postComments {
			postCommentIds[i] = postComment.PostCommentID
		}

		// 대댓글 조회
		replies, err := mysql.PostComments(
			qm.Load(mysql.PostCommentRels.Member),
			qm.WhereIn("parent_post_comment_id IN ?", postCommentIds...), // WhereIn 사용
			qm.WhereNotIn("member_id NOT IN ?", blockedMemberIds...),     // WhereNotIn 대신 적절하게 처리
			qm.Where("is_recomment = TRUE"),                              // 단일 조건은 따로 처리
			qm.OrderBy("post_comment.created_at ASC"),
		).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		repliesIds := make([]interface{}, len(replies))
		for i, reply := range replies {
			repliesIds[i] = reply.PostCommentID
		}

		// 전체 댓글 수를 가져오는 쿼리
		totalCommentsCount, err := mysql.PostComments(
			qm.Where("post_comment.post_id = ? AND post_comment.deleted_at IS NULL", postId),
			qm.WhereNotIn("post_comment.member_id NOT IN ?", blockedMemberIds...), // 블랙리스트 제외
		).Count(c.Request.Context(), db)

		// 결과가 없는 경우 빈 리스트 반환
		if len(postComments) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentV2Response{TotalPostCommentCount: int(totalCommentsCount), PostComments: []PostCommentV2Response{}, LastCursor: int64(cursorInt)}) //마지막(가장 최근 id) 커서값, 없으면 0
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

		likeCountList := append(repliesIds, postCommentIDs...)

		// 해당 song_id와 member_id에 대한 모든 좋아요를 가져오기
		likes, err := mysql.PostCommentLikes(
			qm.WhereIn("post_comment_id IN ?", likeCountList...), // `IN` 조건에 슬라이스 전달
			qm.Where("member_id = ?", blockerId),                 // `AND` 대신 `Where`로 처리
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

		repliesMap := make(map[int64][]PostCommentV2Response)
		for _, reply := range replies {
			repliesMap[reply.ParentPostCommentID.Int64] = append(repliesMap[reply.ParentPostCommentID.Int64], PostCommentV2Response{
				PostCommentId:       reply.PostCommentID,
				Content:             reply.Content.String,
				IsRecomment:         reply.IsRecomment.Bool,
				ParentPostCommentId: reply.ParentPostCommentID.Int64,
				PostId:              reply.PostID,
				MemberId:            reply.MemberID,
				Nickname:            reply.R.Member.Nickname.String,
				CreatedAt:           reply.CreatedAt.Time,
				Likes:               reply.Likes,
				IsLiked:             likedCommentMap[reply.PostCommentID],
			})
		}

		// Initialize a slice to hold all comments
		topLevelComments := make([]PostCommentV2Response, 0, sizeInt)

		// Add all top-level comments (those without parent comments) to the slice
		for _, comment := range postComments {
			if !comment.IsRecomment.Bool {
				if songOnPostMap[comment.PostCommentID] == nil {
					songOnPostMap[comment.PostCommentID] = []SongOnPost{}
				}
				// Top-level comment, add to slice
				topLevelComments = append(topLevelComments, PostCommentV2Response{
					PostCommentId:       comment.PostCommentID,
					Content:             comment.Content.String,
					IsRecomment:         comment.IsRecomment.Bool,
					ParentPostCommentId: comment.ParentPostCommentID.Int64,
					PostId:              comment.PostID,
					MemberId:            comment.MemberID,
					Nickname:            comment.NickName.String,
					CreatedAt:           comment.CreatedAt.Time,
					Likes:               comment.Likes,
					IsLiked:             likedCommentMap[comment.PostCommentID],
					SongOnPostComment:   songOnPostMap[comment.PostCommentID], // Add the song list here
					PostRecommentCount:  comment.ReplyCount,
					PostRecomments:      repliesMap[comment.PostCommentID],
				})
			}
		}

		// 다음 페이지를 위한 커서 값 설정
		lastCursor := int64(cursorInt)
		if len(topLevelComments) > 0 {
			lastCursor = topLevelComments[len(topLevelComments)-1].PostCommentId
		}

		// Return comments as part of the response
		pkg.BaseResponse(c, http.StatusOK, "success", GetPostCommentV2Response{TotalPostCommentCount: int(totalCommentsCount), PostComments: topLevelComments, LastCursor: lastCursor})
	}
}
