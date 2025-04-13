package repository_test

import (
	"SingSong-Server/internal/repository"
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBlockedMemberIDs(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "blocker_member_id", "blocked_member_id", "created_at", "updated_at"}).
		AddRow(1, 100, 1, time.Now(), time.Now()).
		AddRow(2, 100, 2, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `blacklist`.* FROM `blacklist` WHERE (blocker_member_id = ?)")).
		WithArgs(100).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := repository.GetBlockedMemberIDs(ctx, db, 100)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), result[0])
	assert.Equal(t, int64(2), result[1])
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCommentCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

	expectedSQL := regexp.QuoteMeta("SELECT COUNT(*) FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL)")
	mock.ExpectQuery(expectedSQL).
		WithArgs(123).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := repository.GetCommentCount(ctx, db, 123)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTopLevelComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"comment_id", "song_info_id", "member_id", "content", "created_at"}).
		AddRow(1, 123, 456, "test comment", time.Now())

	// 실제 쿼리 패턴에 맞게 정확히 업데이트
	expectedSQL := regexp.QuoteMeta("SELECT `comment`.* FROM `comment` WHERE (comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0) AND (1=1) AND (comment.comment_id > ?) ORDER BY comment.comment_id ASC LIMIT 10;")
	mock.ExpectQuery(expectedSQL).
		WithArgs(123, 0).
		WillReturnRows(rows)

	// Member 테이블 쿼리 추가 (경우에 따라 필요할 수 있음)
	memberSQL := regexp.QuoteMeta("SELECT * FROM `member` WHERE (`member`.`member_id` IN (?))")
	memberRows := sqlmock.NewRows([]string{"member_id", "name", "email"}).
		AddRow(456, "Test User", "test@example.com")
	mock.ExpectQuery(memberSQL).
		WithArgs(456).
		WillReturnRows(memberRows)

	ctx := context.Background()
	result, err := repository.GetTopLevelComments(ctx, db, 123, "comment.comment_id > ?", 0, "comment.comment_id ASC", 10, []interface{}{})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRecomments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 테스트용 데이터
	rows := sqlmock.NewRows([]string{
		"comment_id", "parent_comment_id", "member_id", "content", "created_at",
	}).AddRow(2, 1, 789, "recomment", time.Now())

	// regexp.QuoteMeta 사용
	expectedSQL := regexp.QuoteMeta("SELECT `comment`.* FROM `comment` WHERE (comment.deleted_at IS NULL) AND (`parent_comment_id` IN (?)) AND (`comment`.`member_id` NOT IN (?)) ORDER BY comment.created_at ASC")
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 789).
		WillReturnRows(rows)

	// Member 테이블 쿼리도 mock 처리 추가
	memberSQL := regexp.QuoteMeta("SELECT * FROM `member` WHERE (`member`.`member_id` IN (?))")
	memberRows := sqlmock.NewRows([]string{"member_id", "name", "email"}).
		AddRow(789, "Test User", "test@example.com")
	mock.ExpectQuery(memberSQL).
		WithArgs(789).
		WillReturnRows(memberRows)

	ctx := context.Background()
	result, err := repository.GetRecomments(ctx, db, []interface{}{1}, []interface{}{789})
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLikesForComments(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"comment_id", "member_id", "created_at"}).
		AddRow(1, 100, time.Now()).
		AddRow(2, 100, time.Now())

	// 세미콜론을 포함한 정확한 쿼리 패턴
	expectedSQL := regexp.QuoteMeta("SELECT `comment_like`.* FROM `comment_like` WHERE (`comment_id` IN (?,?)) AND (member_id = ?) AND (deleted_at is null);")
	mock.ExpectQuery(expectedSQL).
		WithArgs(1, 2, 100).
		WillReturnRows(rows)

	ctx := context.Background()
	result, err := repository.GetLikesForComments(ctx, db, []interface{}{1, 2}, 100)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
