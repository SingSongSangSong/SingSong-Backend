package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"context"
	"database/sql"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
)

// todo: https://firebase.google.com/docs/cloud-messaging/send-message?hl=ko#go

// TestNotification godoc
// @Summary      알림 전송 테스트
// @Description  알림 전송 테스트
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v1/notifications/test [post]
func TestNotification(db *sql.DB, firebaseApp *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		client, err := firebaseApp.Messaging(context.Background())
		if err != nil {
			log.Printf("error getting Messaging client: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 안드로이드 테스트
		message := &messaging.Message{
			Token: "invalidtoken", // 알림을 보낼 대상 클라이언트의 FCM 토큰
			Notification: &messaging.Notification{
				Title: "이건 제목이구",
				Body:  "이건 body란다",
			},
			Data: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		}

		// 메시지 전송에 타임아웃을 주고 싶다면 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) 이걸 쓸 수 있다
		_, err = client.Send(context.Background(), message)
		if err != nil {
			log.Printf("error sending android message: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "error - "+err.Error(), nil)
	}
}

// ScreenType 정의
type ScreenType string

// 허용되는 알림 타입
const (
	// 필요한 다른 알림 타입을 추가 가능
	SongScreen ScreenType = "SONG"
	PostScreen ScreenType = "POST"
)

type NotificationMessage struct {
	Title             string  // 알림 제목
	Body              string  // 알림 내용
	SenderMemberId    int64   // 발신자 ID
	ReceiverMemberIds []int64 // 수신자 ID
	ScreenType        ScreenType
	ScreenTypeId      int64 // 게시글ID or 쏭ID
}

func SendNotification(db *sql.DB, firebaseApp *firebase.App, notificationMessage NotificationMessage) {
	ctx := context.Background()
	client, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Printf("error getting Messaging client: %v", err)
		return
	}

	receiverIds := notificationMessage.ReceiverMemberIds
	ids := make([]interface{}, len(receiverIds))
	for i, v := range receiverIds {
		ids[i] = v
	}
	if len(ids) == 0 {
		return
	}

	all, err := mysql.MemberDeviceTokens(
		qm.WhereIn("member_id in ?", ids...),
		qm.Where("is_activate = true"),
	).All(ctx, db)
	if err != nil {
		log.Printf("error - "+err.Error(), err)
		return
	}

	registrationTokens := make([]string, 0, len(all))
	for _, device := range all {
		registrationTokens = append(registrationTokens, device.DeviceToken)
	}

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: notificationMessage.Title,
			Body:  notificationMessage.Body,
		},
		Tokens: registrationTokens,
	}

	br, err := client.SendEachForMulticast(ctx, message)
	if err != nil {
		log.Printf("error sending notifications - " + err.Error())
		return
	}
	if br.FailureCount > 0 {
		var failedTokens []string
		for _, resp := range br.Responses {
			if !resp.Success {
				failedTokens = append(failedTokens, resp.Error.Error())
			}
		}
		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	}
}

func SaveNotificationHistory(db *sql.DB, notificationMessage NotificationMessage) {

}

// 게시글에 댓글이 달렸다는 알림 보내기 => 게시글 작성자에게는 꼭 알림
func NotifyCommentOnPost(db *sql.DB, firebaseApp *firebase.App, memberId int64, postId int64, commentContent string) {
	uniqueMemberIds, err := mysql.PostComments(
		qm.Select("DISTINCT member_id"),
		qm.Where("post_id = ?", postId),
	).All(context.Background(), db)
	if err != nil {
		log.Printf("error fetching unique member ids: %v", err)
		return
	}
	receiverIds := make([]int64, len(uniqueMemberIds))
	for i, v := range uniqueMemberIds {
		receiverIds[i] = v.MemberID
	}
	notification := NotificationMessage{
		Title:             "게시글에 새로운 댓글이 달렸어요",
		Body:              commentContent,
		SenderMemberId:    memberId,
		ReceiverMemberIds: receiverIds,
		ScreenType:        PostScreen,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

// 게시글 대댓글이 달렸을 경우엔 게시글 작성자와, 부모댓글/대댓글 작성자들에게 알림
func NotifyRecommentOnPostComment(db *sql.DB, firebaseApp *firebase.App, memberId int64, postId int64, commentContent string) {

}

// 노래 댓글에 답글이 달렸다는 알림 보내기 => 부모댓글/대댓글 작성자들에게 알림
func NotifyRecommentOnSongComment() {

}
