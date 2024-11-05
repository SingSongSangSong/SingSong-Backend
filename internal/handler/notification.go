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
	"strconv"
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
			Token: "test token",
			Notification: &messaging.Notification{
				Title: "안녕 이건 title이고",
				Body:  "이건 body란다",
			},
			Data: map[string]string{
				"screenType": "POST",
				"screenId":   "30",
			},
		}

		// 메시지 전송에 타임아웃을 주고 싶다면 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) 이걸 쓸 수 있다
		_, err = client.Send(context.Background(), message)
		if err != nil {
			log.Printf("error sending android message: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

// ScreenType 정의
type ScreenType string

// 허용되는 알림 타입
const (
	// 필요한 다른 알림 타입을 추가 가능
	SongScreen ScreenType = "SONG"
	PostScreen ScreenType = "POST"
	HomeScreen ScreenType = "HOME"
)

type NotificationMessage struct {
	Title             string  // 알림 제목
	Body              string  // 알림 내용
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
		Data: map[string]string{
			"screenType": string(notificationMessage.ScreenType),
			"screenId":   strconv.FormatInt(notificationMessage.ScreenTypeId, 10),
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
		//todo: request entity was not found 처리필요?
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

func ToUniqueMemberIds(memberIds []int64) []int64 {
	uniqueMemberIds := make(map[int64]struct{})
	for _, id := range memberIds {
		uniqueMemberIds[id] = struct{}{}
	}
	var result []int64
	for id := range uniqueMemberIds {
		result = append(result, id)
	}
	return result
}

//todo: 댓글 작성자한테는 보내면 안됨!, 제목에 게시글 제목이나 댓글내용을 알려줘야하나?, 예외처리 필요

// 게시글에 그냥 댓글이 달렸다는 알림 => 게시글 작성자에게는 꼭 알림
func NotifyCommentOnPost(db *sql.DB, firebaseApp *firebase.App, postId int64, commentContent string) {
	post, err := mysql.Posts(
		qm.Where("post_id = ?", postId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post: %v", err)
		return
	}
	receiverId := make([]int64, 1)
	receiverId = append(receiverId, post.MemberID)

	notification := NotificationMessage{
		Title:             "게시글에 새로운 댓글이 달렸어요!",
		Body:              commentContent,
		ReceiverMemberIds: receiverId,
		ScreenType:        PostScreen,
		ScreenTypeId:      postId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

// 게시글 대댓글 => 게시글 작성자와, 부모댓글/대댓글 작성자들에게 알림
func NotifyRecommentOnPostComment(db *sql.DB, firebaseApp *firebase.App, parentPostCommentId int64, postId int64, commentContent string) {
	// 게시글 작성자
	post, err := mysql.Posts(
		qm.Where("post_id = ? and deleted_at is null", postId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("no post to send notification", err)
		return
	}

	// 부모댓글 작성자
	parentComment, err := mysql.PostComments(
		qm.Where("post_comment_id = ? and deleted_at is null", parentPostCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching parent comment: "+err.Error(), err)
		return
	}

	// 답댓글 작성자
	babyComments, err := mysql.PostComments(
		qm.Where("parent_post_comment_id = ? and deleted_at is null", parentPostCommentId),
	).All(context.Background(), db)
	if err != nil {
		log.Printf("error fetching baby post comments: %v", err)
		return
	}

	var receiverIds []int64
	for _, v := range babyComments {
		receiverIds = append(receiverIds, v.MemberID)
	}
	receiverIds = append(receiverIds, post.MemberID)
	receiverIds = append(receiverIds, parentComment.MemberID)

	receiverIds = ToUniqueMemberIds(receiverIds)

	notification := NotificationMessage{
		Title:             "댓글에 새로운 답글이 달렸어요!",
		Body:              commentContent,
		ReceiverMemberIds: receiverIds,
		ScreenType:        PostScreen,
		ScreenTypeId:      postId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

// 노래 댓글에 답글이 달렸다는 알림 보내기 => 부모댓글/대댓글 작성자들에게 알림
func NotifyRecommentOnSongComment(db *sql.DB, firebaseApp *firebase.App, parentCommentId int64, songId int64, commentContent string) {
	// 부모댓글 작성자
	parentComment, err := mysql.Comments(
		qm.Where("comment_id = ? and deleted_at is null", parentCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching parent comment: "+err.Error(), err)
		return
	}

	// 답댓글 작성자
	babyComments, err := mysql.Comments(
		qm.Where("parent_comment_id = ? and deleted_at is null", parentCommentId),
	).All(context.Background(), db)
	if err != nil {
		log.Printf("error baby comments: %v", err)
		return
	}

	var receiverIds []int64
	for _, v := range babyComments {
		receiverIds = append(receiverIds, v.MemberID)
	}
	receiverIds = append(receiverIds, parentComment.MemberID)

	receiverIds = ToUniqueMemberIds(receiverIds)

	notification := NotificationMessage{
		Title:             "댓글에 새로운 답글이 달렸어요!",
		Body:              commentContent,
		ReceiverMemberIds: receiverIds,
		ScreenType:        SongScreen,
		ScreenTypeId:      songId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

//todo:  게시즐 좋아요/ 게시글댓글 좋아요 / 노래댓글 좋아요 알림
