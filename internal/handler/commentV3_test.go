package handler

//
//import (
//	"encoding/json"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"regexp"
//	"testing"
//	"time"
//)
//
//// a successful case
//func TestGetCommentsOnSongV3(t *testing.T) {
//	// 정규식 매칭을 사용하는 sqlmock 생성
//	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
//	assert.NoError(t, err)
//	defer db.Close()
//
//	// Gin 테스트 컨텍스트 설정
//	gin.SetMode(gin.TestMode)
//	w := httptest.NewRecorder()
//	c, _ := gin.CreateTestContext(w)
//
//	// 테스트 요청 설정
//	c.Params = []gin.Param{
//		{Key: "songId", Value: "123"},
//	}
//	c.Request, _ = http.NewRequest("GET", "/api/v3/songs/123/comments?filter=old&size=10&cursor=0", nil)
//	c.Set("memberId", int64(100)) // 현재 사용자 ID 설정
//
//	// 1. GetBlockedMemberIDs 쿼리 설정
//	blacklistRows := sqlmock.NewRows([]string{"id", "blocker_member_id", "blocked_member_id"}).
//		AddRow(1, 100, 789)
//	blacklistSQL := regexp.QuoteMeta("SELECT `blacklist`.* FROM `blacklist` WHERE (blocker_member_id = ?)")
//	mock.ExpectQuery(blacklistSQL).
//		WithArgs(int64(100)).
//		WillReturnRows(blacklistRows)
//
//	// 2. GetTopLevelComments 쿼리 설정
//	commentTime := time.Now()
//	commentRows := sqlmock.NewRows([]string{
//		"comment_id", "song_info_id", "member_id", "content", "is_recomment",
//		"parent_comment_id", "created_at", "likes",
//	}).
//		AddRow(1, 123, 456, "Top level comment", false, 0, commentTime, 5)
//
//	// 실제 로그에서 관찰된 쿼리 그대로 사용
//	commentSQL := regexp.QuoteMeta("SELECT `comment`.* FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0) AND (`comment`.`member_id` NOT IN (?)) AND (comment.comment_id > ?) ORDER BY comment.comment_id ASC LIMIT 10;")
//	mock.ExpectQuery(commentSQL).
//		WithArgs(123, 789, int64(0)).
//		WillReturnRows(commentRows)
//
//	// Member 테이블 eager loading
//	memberRows := sqlmock.NewRows([]string{"member_id", "nickname", "email"}).
//		AddRow(456, "User1", "user1@example.com")
//	memberSQL := regexp.QuoteMeta("SELECT * FROM `member` WHERE (`member`.`member_id` IN (?));")
//	mock.ExpectQuery(memberSQL).
//		WithArgs(456).
//		WillReturnRows(memberRows)
//
//	// 3. GetCommentCount 쿼리 설정
//	countRows := sqlmock.NewRows([]string{"count"}).
//		AddRow(10)
//
//	// 실제 로그에서 관찰된 쿼리 그대로 사용
//	countSQL := regexp.QuoteMeta("SELECT COUNT(*) FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL);")
//	mock.ExpectQuery(countSQL).
//		WithArgs(123).
//		WillReturnRows(countRows)
//
//	// 4. GetRecomments 쿼리 설정
//	recommentTime := time.Now()
//	recommentRows := sqlmock.NewRows([]string{
//		"comment_id", "song_info_id", "member_id", "content", "is_recomment",
//		"parent_comment_id", "created_at", "likes",
//	}).
//		AddRow(2, 123, 567, "Recomment", true, 1, recommentTime, 3)
//
//	// 실제 로그에서 관찰된 쿼리 그대로 사용
//	recommentSQL := regexp.QuoteMeta("SELECT `comment`.* FROM `comment` WHERE (comment.deleted_at IS NULL) AND (`parent_comment_id` IN (?)) AND (`comment`.`member_id` NOT IN (?)) ORDER BY comment.created_at ASC;")
//	mock.ExpectQuery(recommentSQL).
//		WithArgs(1, 789).
//		WillReturnRows(recommentRows)
//
//	// Member 테이블 eager loading (recomment)
//	recommentMemberRows := sqlmock.NewRows([]string{"member_id", "nickname", "email"}).
//		AddRow(567, "User2", "user2@example.com")
//	recommentMemberSQL := regexp.QuoteMeta("SELECT * FROM `member` WHERE (`member`.`member_id` IN (?));")
//	mock.ExpectQuery(recommentMemberSQL).
//		WithArgs(567).
//		WillReturnRows(recommentMemberRows)
//
//	// 5. GetLikesForComments 쿼리 설정
//	likeRows := sqlmock.NewRows([]string{"comment_id", "member_id", "created_at"}).
//		AddRow(1, 100, time.Now())
//
//	// 실제 로그에서 관찰된 쿼리 그대로 사용
//	likeSQL := regexp.QuoteMeta("SELECT `comment_like`.* FROM `comment_like` WHERE (`comment_id` IN (?,?)) AND (member_id = ?) AND (deleted_at is null);")
//	mock.ExpectQuery(likeSQL).
//		WithArgs(1, 2, int64(100)).
//		WillReturnRows(likeRows)
//
//	// 핸들러 함수 호출
//	handler := GetCommentsOnSongV3(db)
//	handler(c)
//
//	// 응답 확인
//	assert.Equal(t, http.StatusOK, w.Code)
//
//	// JSON 응답 파싱
//	var response struct {
//		Code    int                   `json:"code"`
//		Message string                `json:"message"`
//		Data    CommentPageV3Response `json:"data"`
//	}
//	err = json.Unmarshal(w.Body.Bytes(), &response)
//	assert.NoError(t, err)
//
//	// 응답 검증 - 빈 결과 처리를 위한 안전 장치 추가
//	assert.Equal(t, "success", response.Message)
//	assert.Equal(t, int64(10), response.Data.CommentCount)
//	assert.Equal(t, 1, len(response.Data.Comments))
//
//	// 배열이 비어 있지 않은 경우에만 접근
//	if len(response.Data.Comments) > 0 {
//		assert.Equal(t, int64(1), response.Data.Comments[0].CommentId)
//		assert.Equal(t, true, response.Data.Comments[0].IsLiked)
//
//		if len(response.Data.Comments[0].Recomments) > 0 {
//			assert.Equal(t, int64(2), response.Data.Comments[0].Recomments[0].CommentId)
//		}
//	}
//
//	// 모든 mock 기대가 충족되었는지 확인
//	assert.NoError(t, mock.ExpectationsWereMet())
//}
