package user

import (
	"SingSong-Backend/internal/model"
	"SingSong-Backend/internal/pkg"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RedisHandler struct {
	redisModel *model.RedisModel
}

func NewRedisHandler(redisModel *model.RedisModel) (*RedisHandler, error) {
	redisHandler := &RedisHandler{redisModel: redisModel}
	return redisHandler, nil
}

type PublicKeyDto struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

// 공개키 목록 조회 URL 요청 함수
func (handler *RedisHandler) SetPublicKeys(c *gin.Context) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(REQUEST_URL)
	if err != nil {
		log.Printf("HTTP 요청 오류: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		log.Printf("HTTP 응답 오류: 상태 코드 %d", resp.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("응답 본문 읽기 오류: %v", err)
		return
	}

	// PublicKeyDto 구조체 생성
	publicKey := PublicKeyDto{
		Provider: PROVIDER,
		Key:      string(body),
	}

	// JSON으로 변환
	jsonData, err := json.Marshal(publicKey)
	if err != nil {
		log.Printf("JSON 마샬 오류: %v", err)
		return
	}

	// Redis에 저장
	key := handler.redisModel.SavePublicKey(c, PROVIDER, jsonData, 0)
	log.Println("데이터가 성공적으로 Redis에 저장되었습니다." + key.Val())

	pkg.BaseResponse(c, http.StatusOK, "공개키 저장 성공", key)
}

// Redis에서 공개키 가져오기 함수
func (handler *RedisHandler) GetPublicKeys(c *gin.Context, provider string) ([]JsonWebKey, error) {
	response := handler.redisModel.Get(c, provider)
	if err := response.Err(); err != nil {
		log.Printf("오류 발생: %v", err)

		// 공개키 설정 함수 호출
		handler.SetPublicKeys(c)

		// 다시 시도하여 공개키 가져오기
		response = handler.redisModel.Get(c, provider)
		if err := response.Err(); err != nil {
			log.Printf("오류 발생 (다시 시도 후): %v", err)
			return nil, err
		}
	}
	log.Printf("응답: %v", response.Val())
	publicKeyDto, err := parsePublicKeyDto(response.Val())
	if err != nil {
		log.Printf("오류 발생: %v", err)
		return nil, err
	}

	// PublicKeyDto의 key 필드를 파싱하여 각 Key 구조체로 변환
	keys, err := parseKeysFromPublicKeyDto(publicKeyDto)
	if err != nil {
		log.Printf("오류 발생: %v", err)
		return nil, err
	}

	return keys, nil
}

// JSON 데이터를 파싱하여 PublicKeyDto 구조체로 변환
func parsePublicKeyDto(jsonData string) (*PublicKeyDto, error) {
	var publicKeyDto PublicKeyDto
	err := json.Unmarshal([]byte(jsonData), &publicKeyDto)
	if err != nil {
		return nil, fmt.Errorf("JSON 언마샬 오류: %v", err)
	}
	return &publicKeyDto, nil
}

// PublicKeyDto의 key 필드를 파싱하여 각 Key 구조체로 변환
func parseKeysFromPublicKeyDto(publicKeyDto *PublicKeyDto) ([]JsonWebKey, error) {
	var keyContainer KeyContainer
	err := json.Unmarshal([]byte(publicKeyDto.Key), &keyContainer)
	if err != nil {
		return nil, fmt.Errorf("키 필드 JSON 언마샬 오류: %v", err)
	}
	return keyContainer.Keys, nil
}
