package handler

import (
	"SingSong-Server/internal/db/mysql"
	"context"
	"database/sql"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"time"
)

func ActivateDeviceToken(db *sql.DB, deviceToken string, memberId int64) {
	if deviceToken == "" {
		return
	}
	ctx := context.Background()
	exists, err := mysql.MemberDeviceTokens(
		qm.Where("member_id = ?", memberId),
		qm.And("device_token = ?", deviceToken),
	).Exists(ctx, db)
	if err != nil {
		log.Printf("error -" + err.Error())
		return
	}
	if exists {
		// activate하게 만들기
		_, err := mysql.MemberDeviceTokens(
			qm.Where("member_id = ?", memberId),
			qm.And("device_token = ?", deviceToken),
		).UpdateAll(ctx, db, mysql.M{
			"is_activate": null.BoolFrom(true),
		})
		if err != nil {
			log.Printf("error -" + err.Error())
			return
		}
	} else {
		token := mysql.MemberDeviceToken{MemberID: memberId, DeviceToken: deviceToken, IsActivate: null.BoolFrom(true)}
		err := token.Insert(ctx, db, boil.Infer())
		if err != nil {
			log.Printf("error -" + err.Error())
			return
		}
	}
}

func InvalidateAllDeviceTokens(db *sql.DB, memberId int64) {
	ctx := context.Background()
	rowsAffected, err := mysql.MemberDeviceTokens(
		qm.Where("member_id = ?", memberId),
	).UpdateAll(ctx, db, mysql.M{
		"is_activate": null.BoolFrom(false),
		"deleted_at":  null.TimeFrom(time.Now()),
	})
	if err != nil {
		log.Printf("error -" + err.Error())
	}
	log.Printf("Success! Rows affected: %d", rowsAffected)
}
