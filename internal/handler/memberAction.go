package handler

import (
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"strconv"
)

func logMemberAction(db *sql.DB, memberId interface{}, actionType string, actionScore float32, songInfoIds ...string) {
	memberIdInt64, ok := memberId.(int64)
	if !ok {
		log.Printf("memberId is not of type int64")
		return
	}

	ctx := context.Background()
	one, err := mysql.Members(qm.Where("member_id = ? AND deleted_at is null", memberIdInt64)).One(ctx, db)
	if err != nil {
		log.Printf("failed to get member: %s", err.Error())
		return
	}

	for _, songInfoIdStr := range songInfoIds {
		songInfoIdInt64, err := strconv.ParseInt(songInfoIdStr, 10, 64)
		if err != nil {
			log.Printf("failed to parse songInfoId: %s", err.Error())
			continue
		}

		action := mysql.MemberAction{
			MemberID:    memberIdInt64,
			Gender:      one.Gender,
			Birthyear:   one.Birthyear,
			SongInfoID:  songInfoIdInt64,
			ActionType:  actionType,
			ActionScore: actionScore,
		}

		if err := action.Insert(ctx, db, boil.Infer()); err != nil {
			log.Printf("failed to insert member action: %s", err.Error())
		}
	}
}
