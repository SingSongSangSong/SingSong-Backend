package pkg

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// sentry에서 에러를 봐야 할 부분에만 적용하기(ex. 500 에러)
func SendToSentryWithStack(c *gin.Context, err error) {
	//stack := debug.Stack()
	hub := sentrygin.GetHubFromContext(c)
	if hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			//scope.SetExtra("stack_trace", string(stack)) // 스택 트레이스 추가
			hub.CaptureException(err)
		})
	}
}
