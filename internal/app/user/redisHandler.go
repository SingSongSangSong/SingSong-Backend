package user

import (
	"SingSong-Backend/internal/model"
	"SingSong-Backend/internal/pkg"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type RedisHandler struct {
	redisModel *model.RedisModel
}

func NewRedisHandler(redisModel *model.RedisModel) (*RedisHandler, error) {
	redisHandler := &RedisHandler{redisModel: redisModel}
	return redisHandler, nil
}

const (
	REQUEST_URL = "https://kauth.kakao.com/.well-known/jwks.json" // 여기에 실제 요청 URL을 입력하세요.
	PROVIDER    = "KAKAO"                                         // 여기에 실제 공급자 이름을 입력하세요.
)

// PublicKey 구조체 정의
type PublicKey struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

func (handler *RedisHandler) SetPublicKeys(ctx context.Context) {
	// 공개키 목록 조회 URL 요청
	resp, err := http.Get(REQUEST_URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 응답 코드 확인
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		log.Printf("HTTP 요청 실패: %d", resp.StatusCode)
	}

	// 응답 본문 읽기
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("응답 본문 읽기 실패: %v", err)
	}

	// JSON 데이터 저장
	publicKey := PublicKey{
		Provider: PROVIDER,
		Key:      string(body),
	}

	jsonData, err := json.Marshal(publicKey)
	if err != nil {
		pkg.BaseResponse(ctx, http.StatusInternalServerError, "JSON 데이터 변환 실패", nil)
	}

	// Redis에 저장
	statusCmd := handler.redisModel.Set(ctx, "public_key", jsonData)
	if err := statusCmd.Err(); err != nil {
		return
	}

	return
}
