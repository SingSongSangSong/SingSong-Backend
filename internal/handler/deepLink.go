package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"
	"fmt"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/url"
	"time"
)

func CreateSongDeepLink(db *sql.DB, songInfoId int64) (string, error) {
	one, err := mysql.SongInfos(qm.Where("song_info_id = ?", songInfoId)).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching song info: %v", err)
		return "", err
	}

	// 문자열 생성
	result := fmt.Sprintf(
		"%ssong?songId=%d&&songName=%s&&singerName=%s&&album=%s&&melonLink=%s&&isMr=%t&&isLive=%t",
		conf.NotificationConfigInstance.DeepLinkBase,
		one.SongInfoID,
		one.SongName,
		one.ArtistName,
		url.QueryEscape(one.Album.String),
		url.QueryEscape(CreateMelonLinkByMelonSongId(one.MelonSongID)),
		one.IsMR.Bool,
		one.IsLive.Bool,
	)

	return result, nil
}

func CreatePostDeepLink(db *sql.DB, postId int64) (string, error) {
	one, err := mysql.Posts(
		qm.Load(mysql.PostRels.PostComments),
		qm.Load(mysql.PostRels.Member),
		qm.LeftOuterJoin("post_comment on post_comment.post_id = post.post_id"),
		qm.Where("post.post_id = ?", postId),
		qm.Where("post.deleted_at is null"),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post: %v", err)
		return "", err
	}

	//deleted_at이 NULL인 댓글만 카운트
	validCommentCount := 0
	for _, comment := range one.R.PostComments {
		if !comment.DeletedAt.Valid { // deleted_at이 NULL인지 확인
			validCommentCount++
		}
	}

	// 문자열 생성
	result := fmt.Sprintf(
		"%splayground?postId=%d&&title=%s&&content=%s&&createdAt=%s&&nickname=%s&&likes=%d&&commentCount=%d",
		conf.NotificationConfigInstance.DeepLinkBase,
		one.PostID,
		url.QueryEscape(one.Title),
		url.QueryEscape(one.Content.String),
		url.QueryEscape(one.CreatedAt.Time.Format(time.RFC3339)),
		url.QueryEscape(one.R.Member.Nickname.String),
		one.Likes,
		validCommentCount,
	)

	return result, nil
}

func CreateHomeDeepLink() string {
	return conf.NotificationConfigInstance.DeepLinkBase + "home"
}
