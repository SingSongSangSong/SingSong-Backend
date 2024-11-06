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
	"strings"
)

// todo: https://firebase.google.com/docs/cloud-messaging/send-message?hl=ko#go
// todo: 댓글 작성자/게시글 작성자한테는 알림을 보내면 안됨(테스트 필요) + 그밖에 예외처리 필요(중요!!) (+ 제목에 게시글 제목이나 댓글내용을 알려줘야하는지 생각해보기!)
// todo: 내 알림 조회 api

type AnnouncementRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendAnnouncementNotification godoc
// @Summary      디바이스 토큰이 활성화된 모든 유저에게 공지사항 전송
// @Description  공지사항 전송
// @Tags         Notification
// @Accept       json
// @Produce      json
// @Param        AnnouncementRequest  body   AnnouncementRequest  true  "알림 내용"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /v1/notifications/announcements [post]
func SendAnnouncement(db *sql.DB, firebaseApp *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		announcementRequest := &AnnouncementRequest{}
		if err := c.ShouldBindJSON(&announcementRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		members, err := mysql.Members(
			qm.Where("deleted_at is null"),
		).All(context.Background(), db)
		if err != nil {
			log.Printf("error fetching members: %v", err)
			return
		}
		receiverIds := make([]int64, len(members))
		for i, v := range members {
			receiverIds[i] = v.MemberID
		}

		notification := NotificationMessage{
			Title:             announcementRequest.Title,
			Body:              announcementRequest.Body,
			ReceiverMemberIds: receiverIds,
			ScreenType:        HomeScreen,
			ScreenTypeId:      0,
		}
		go SendNotification(db, firebaseApp, notification)
		go SaveNotificationHistory(db, notification)

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
		qm.Where("is_activate = true and deleted_at is null"),
	).All(ctx, db)
	if err != nil {
		log.Printf("error - "+err.Error(), err)
		return
	}

	registrationTokens := make([]string, 0, len(all))
	for _, device := range all {
		registrationTokens = append(registrationTokens, device.DeviceToken)
	}

	// 500개씩 나누어서 전송
	const batchSize = 500
	for start := 0; start < len(registrationTokens); start += batchSize {
		end := start + batchSize
		if end > len(registrationTokens) {
			end = len(registrationTokens)
		}
		batchTokens := registrationTokens[start:end]

		//todo: 딥링크 추가
		message := &messaging.MulticastMessage{
			Notification: &messaging.Notification{
				Title: notificationMessage.Title,
				Body:  notificationMessage.Body,
			},
			Data: map[string]string{
				"screenType": string(notificationMessage.ScreenType),
				"screenId":   strconv.FormatInt(notificationMessage.ScreenTypeId, 10),
			},
			Android: &messaging.AndroidConfig{
				Priority: "high",
				Notification: &messaging.AndroidNotification{
					Sound:     "default",
					ChannelID: "default-channel-id",
				},
			},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{
					Aps: &messaging.Aps{
						Sound: "default",
					},
				},
			},
			Tokens: batchTokens,
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
}

// 알림 허용을 하지 않은 사람도 앱 내에서 알림을 볼 수 있도록 하기 위해 따로 저장함
func SaveNotificationHistory(db *sql.DB, notificationMessage NotificationMessage) {
	ctx := context.Background()
	if len(notificationMessage.ReceiverMemberIds) == 0 { // todo: 본인에게 보낸 알림을 제외하긴 해야 됨
		return
	}

	// 배치 삽입
	values := make([]interface{}, 0, len(notificationMessage.ReceiverMemberIds)*5)
	placeholders := make([]string, len(notificationMessage.ReceiverMemberIds))

	for i, memberId := range notificationMessage.ReceiverMemberIds {
		values = append(values,
			memberId,
			notificationMessage.Title,
			notificationMessage.Body,
			string(notificationMessage.ScreenType),
			notificationMessage.ScreenTypeId,
		)
		placeholders[i] = "(?, ?, ?, ?, ?, false)"
	}

	query := fmt.Sprintf(`
		INSERT INTO notification_history (member_id, title, body, screen_type, screen_type_id, is_read) 
		VALUES %s`, strings.Join(placeholders, ","))

	_, err := db.ExecContext(ctx, query, values...)
	if err != nil {
		log.Printf("error inserting notification history: %v", err)
	}
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

// 게시글에 그냥 댓글이 달렸다는 알림 => 게시글 작성자에게는 꼭 알림
func NotifyCommentOnPost(db *sql.DB, firebaseApp *firebase.App, senderId int64, postId int64, commentContent string) {
	post, err := mysql.Posts(
		qm.Where("post_id = ? and deleted_at is null", postId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post: %v", err)
		return
	}
	if post.MemberID == senderId {
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
func NotifyRecommentOnPostComment(db *sql.DB, firebaseApp *firebase.App, senderId int64, parentPostCommentId int64, postId int64, commentContent string) {
	// 게시글 작성자
	post, err := mysql.Posts(
		qm.Where("post_id = ? and deleted_at is null", postId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post: %v", err)
		return
	}
	if post.MemberID == senderId {
		return
	}

	// 부모댓글 작성자
	parentComment, err := mysql.PostComments(
		qm.Where("post_comment_id = ? and deleted_at is null", parentPostCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching parent post comment: "+err.Error(), err)
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
		if v.MemberID != senderId { // senderId와 다른 멤버만 추가
			receiverIds = append(receiverIds, v.MemberID)
		}
	}
	if post.MemberID != senderId {
		receiverIds = append(receiverIds, post.MemberID)
	}
	if parentComment.MemberID != senderId {
		receiverIds = append(receiverIds, parentComment.MemberID)
	}

	if len(receiverIds) == 0 || receiverIds == nil { // 알림 보낼게 없다면 리턴
		return
	}

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
func NotifyRecommentOnSongComment(db *sql.DB, firebaseApp *firebase.App, senderId int64, parentCommentId int64, songId int64, commentContent string) {
	// 부모댓글 작성자
	parentComment, err := mysql.Comments(
		qm.Where("comment_id = ? and deleted_at is null", parentCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching parent song comment: "+err.Error(), err)
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
		if v.MemberID != senderId {
			receiverIds = append(receiverIds, v.MemberID)
		}
	}
	if parentComment.MemberID != senderId {
		receiverIds = append(receiverIds, parentComment.MemberID)
	}

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

// 게시글 좋아요 -> 게시글 작성자에게 알림
func NotifyLikeOnPost(db *sql.DB, firebaseApp *firebase.App, senderId int64, postId int64, postTitle string) {
	post, err := mysql.Posts(
		qm.Where("post_id = ? and deleted_at is null", postId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post: %v", err)
		return
	}
	if post.MemberID == senderId {
		return
	}

	receiverId := make([]int64, 1)
	receiverId = append(receiverId, post.MemberID)

	notification := NotificationMessage{
		Title:             "누군가 당신의 게시글에 좋아요를 눌렀어요!",
		Body:              postTitle,
		ReceiverMemberIds: receiverId,
		ScreenType:        PostScreen,
		ScreenTypeId:      postId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

// 게시글 댓글 좋아요 -> 댓글 작성자에게 알림
func NotifyLikeOnPostComment(db *sql.DB, firebaseApp *firebase.App, senderId int64, postCommentId int64, postId int64, commentContent string) {
	postComment, err := mysql.PostComments(
		qm.Where("post_comment_id = ? and deleted_at is null", postCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching post comment: "+err.Error(), err)
		return
	}
	if postComment.MemberID == senderId {
		return
	}

	receiverId := make([]int64, 1)
	receiverId = append(receiverId, postComment.MemberID)

	notification := NotificationMessage{
		Title:             "누군가 당신의 댓글에 좋아요를 눌렀어요!",
		Body:              commentContent,
		ReceiverMemberIds: receiverId,
		ScreenType:        PostScreen,
		ScreenTypeId:      postId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}

// 노래 댓글 좋아요 -> 노래 댓글 작성자에게 알림
func NotifyLikeOnSongComment(db *sql.DB, firebaseApp *firebase.App, senderId int64, songCommentId int64, songId int64, commentContent string) {
	songComment, err := mysql.Comments(
		qm.Where("comment_id = ? and deleted_at is null", songCommentId),
	).One(context.Background(), db)
	if err != nil {
		log.Printf("error fetching song comment: "+err.Error(), err)
		return
	}
	if songComment.MemberID == senderId {
		return
	}

	receiverId := make([]int64, 1)
	receiverId = append(receiverId, songComment.MemberID)

	notification := NotificationMessage{
		Title:             "누군가 당신의 댓글에 좋아요를 눌렀어요!",
		Body:              commentContent,
		ReceiverMemberIds: receiverId,
		ScreenType:        SongScreen,
		ScreenTypeId:      songId,
	}
	SendNotification(db, firebaseApp, notification)
	SaveNotificationHistory(db, notification)
}
