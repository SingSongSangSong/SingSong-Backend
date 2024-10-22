package handler

import (
	"SingSong-Server/internal/pkg"
	"database/sql"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

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
		// 클라이언트 초기화
		client, err := firebaseApp.Messaging(c)
		if err != nil {
			log.Printf("error getting Messaging client: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// 메시지 생성
		message := &messaging.Message{
			Token: "token", // 알림을 보낼 대상 클라이언트의 FCM 토큰
			Notification: &messaging.Notification{
				Title: "이건 알림 제목이양",
				Body:  "이거는 바디 내용이양",
			},
		}

		// 메시지 전송에 타임아웃을 주고 싶다면 ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) 이걸 쓸 수 있다
		_, err = client.Send(c, message)
		if err != nil {
			log.Printf("error sending message: %v", err)
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "error - "+err.Error(), nil)
	}
}
