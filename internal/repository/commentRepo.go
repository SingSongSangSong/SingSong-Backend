package repository

import (
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetBlockedMemberIDs(ctx context.Context, db *sql.DB, blockerId interface{}) ([]interface{}, error) {
	blacklists, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(ctx, db)
	if err != nil {
		return nil, err
	}
	ids := make([]interface{}, 0, len(blacklists))
	for _, b := range blacklists {
		ids = append(ids, b.BlockedMemberID)
	}
	return ids, nil
}

func GetTopLevelComments(ctx context.Context, db *sql.DB, songId int, cursorCondition string, cursorInt int64, orderBy string, size int, blockedIDs []interface{}) (mysql.CommentSlice, error) {
	return mysql.Comments(
		qm.Load(mysql.CommentRels.Member),
		qm.Where("comment.song_info_id = ? AND comment.deleted_at IS NULL AND parent_comment_id = 0", songId),
		qm.WhereNotIn("comment.member_id NOT IN ?", blockedIDs...),
		qm.And(cursorCondition, cursorInt),
		qm.OrderBy(orderBy),
		qm.Limit(size),
	).All(ctx, db)
}

func GetCommentCount(ctx context.Context, db *sql.DB, songId int) (int64, error) {
	return mysql.Comments(
		qm.Where("comment.song_info_id = ? AND comment.deleted_at IS NULL", songId),
	).Count(ctx, db)
}

func GetRecomments(ctx context.Context, db *sql.DB, commentIDs []interface{}, blockedIDs []interface{}) (mysql.CommentSlice, error) {
	return mysql.Comments(
		qm.Load(mysql.CommentRels.Member),
		qm.Where("comment.deleted_at IS NULL"),
		qm.WhereIn("parent_comment_id IN ?", commentIDs...),
		qm.WhereNotIn("comment.member_id NOT IN ?", blockedIDs...),
		qm.OrderBy("comment.created_at ASC"),
	).All(ctx, db)
}

func GetLikesForComments(ctx context.Context, db *sql.DB, commentIDs []interface{}, memberId interface{}) (mysql.CommentLikeSlice, error) {
	return mysql.CommentLikes(
		qm.WhereIn("comment_id IN ?", commentIDs...),
		qm.And("member_id = ?", memberId),
		qm.And("deleted_at is null"),
	).All(ctx, db)
}
