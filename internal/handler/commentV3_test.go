package handler

//import (
//	"SingSong-Server/internal/pkg"
//	"database/sql"
//	"fmt"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"net/http"
//	"regexp"
//	"testing"
//	"time"
//)
//
//func TestGetCommentsOnSongV3(t *testing.T) {
//	db, mock, err := pkg.NewTestDBMock()
//	require.NoError(t, err)
//	defer db.Close()
//
//	handler := GetCommentsOnSongV3(db)
//
//	tests := []struct {
//		name     string
//		filter   string
//		cursor   string
//		cursorOp string
//		orderBy  string
//	}{
//		{"oldest_first", "old", "0", ">", "ASC"},
//		{"newest_first", "recent", "9223372036854775807", "<", "DESC"},
//	}
//
//	for _, tc := range tests {
//		t.Run(tc.name, func(t *testing.T) {
//			c, w := pkg.NewTestGinContext("GET", "/v3/songs/1/comments", "")
//			c.Params = gin.Params{gin.Param{Key: "songId", Value: "1"}}
//			c.Request.URL.RawQuery = fmt.Sprintf("filter=%s&cursor=%s&size=10", tc.filter, tc.cursor)
//			c.Set("memberId", int64(42))
//
//			// 1. 실제로 첫 번째로 실행되는 쿼리는 blacklist 조회입니다
//			mock.ExpectQuery(regexp.QuoteMeta(
//				"SELECT `blacklist`.* FROM `blacklist` WHERE (blocker_member_id = ?)")).
//				WithArgs(int64(42)).
//				WillReturnRows(sqlmock.NewRows([]string{"blacklist_id", "blocker_member_id", "blocked_member_id"}))
//
//			// 2. 다음으로 comment count 쿼리가 실행됩니다
//			mock.ExpectQuery(regexp.QuoteMeta(
//				"SELECT COUNT(*) FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL)")).
//				WithArgs(int64(1)).
//				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
//
//			// 3. 댓글 목록 조회 - 실제 실행되는 쿼리 기준
//			var commentQuery string
//			if tc.filter == "old" {
//				commentQuery = "SELECT `comment`.* FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0) AND (comment.comment_id > ?) ORDER BY comment.comment_id ASC LIMIT ?"
//			} else {
//				commentQuery = "SELECT `comment`.* FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0) AND (comment.comment_id < ?) ORDER BY comment.comment_id DESC LIMIT ?"
//			}
//
//			now := time.Now()
//			// 실제 반환되는 테이블 구조에 맞게 컬럼 설정
//			mock.ExpectQuery(regexp.QuoteMeta(commentQuery)).
//				WithArgs(int64(1), sqlmock.AnyArg(), 10).
//				WillReturnRows(sqlmock.NewRows([]string{
//					"comment_id", "content", "is_recomment", "parent_comment_id",
//					"song_info_id", "member_id", "created_at", "updated_at", "deleted_at", "likes",
//				}).AddRow(
//					int64(1), sql.NullString{String: "Hello world", Valid: true},
//					sql.NullBool{Bool: false, Valid: true},
//					sql.NullInt64{Int64: 0, Valid: false},
//					int64(1), int64(42),
//					now, now, nil, sql.NullInt32{Int32: 5, Valid: true},
//				))
//
//			// 3.1 Member 릴레이션 로드
//			mock.ExpectQuery(regexp.QuoteMeta(
//				"SELECT `member`.* FROM `member` WHERE (`member`.member_id IN (?))")).
//				WithArgs(int64(42)).
//				WillReturnRows(sqlmock.NewRows([]string{
//					"member_id", "nickname", "profile_image_url",
//				}).AddRow(
//					int64(42), sql.NullString{String: "TestUser", Valid: true}, sql.NullString{String: "", Valid: false},
//				))
//
//			// 4. 대댓글 조회
//			mock.ExpectQuery(regexp.QuoteMeta(
//				"SELECT `comment`.* FROM `comment` WHERE (comment.deleted_at IS NULL) AND (parent_comment_id IN (?))")).
//				WithArgs(int64(1)).
//				WillReturnRows(sqlmock.NewRows([]string{
//					"comment_id", "content", "is_recomment", "parent_comment_id",
//					"song_info_id", "member_id", "created_at", "updated_at", "deleted_at", "likes",
//				}))
//
//			// 5. 좋아요 조회
//			mock.ExpectQuery(regexp.QuoteMeta(
//				"SELECT `comment_likes`.* FROM `comment_likes` WHERE (comment_id IN (?) AND member_id = ? AND deleted_at is null)")).
//				WithArgs(int64(1), int64(42)).
//				WillReturnRows(sqlmock.NewRows([]string{
//					"comment_like_id", "comment_id", "member_id", "created_at", "updated_at", "deleted_at",
//				}).AddRow(1, 1, 42, now, now, nil))
//
//			handler(c)
//
//			assert.Equal(t, http.StatusOK, w.Code)
//			assert.Contains(t, w.Body.String(), "\"commentCount\":1")
//			assert.Contains(t, w.Body.String(), "\"content\":\"Hello world\"")
//			assert.Contains(t, w.Body.String(), "\"isLiked\":true")
//		})
//	}
//}
