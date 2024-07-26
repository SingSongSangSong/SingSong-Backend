package handler

import (
	"SingSong-Server/internal/pkg"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	REQUEST_URL = "https://kauth.kakao.com/.well-known/jwks.json" // 공개키 목록 조회 URL
	PROVIDER    = "KAKAO"                                         // 공급자
)

type PublicKeyDto struct {
	Provider string `json:"provider"`
	Key      string `json:"key"`
}

// JsonWebKey struct definition
type JsonWebKey struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Claims 구조체 정의 (필요에 따라 조정 가능)
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// 각 Key 구조체를 담을 구조체 정의
type KeyContainer struct {
	Keys []JsonWebKey `json:"keys"`
}

type LoginRequest struct {
	IdToken  string `json:"IdToken"`
	Provider string `json:"Provider"`
}

func OAuth(c *gin.Context, redis *redis.Client, db *sql.DB) {
	loginRequest := &LoginRequest{}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}

	// email 및 nickname 추출
	payload, err := GetUserEmailFromIdToken(c, redis, loginRequest.IdToken)
	if err != nil {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}

	// email이 없을 경우 에러 반환
	if payload.Email == "" {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - Email is empty", nil)
		return
	}
	// nickname이 없을 경우 에러 반환
	//if(payload.Nickname == "") {
	//	pkg.BaseResponse(c, http.StatusBadRequest, "error - Nickname is empty", nil)
	//	return
	//}

	// email이 db에 있는지 확인
	rows, err := db.Query("SELECT * FROM user WHERE email = ? and provider = ?", payload.Email, loginRequest.Provider)
	if err != nil {
		log.Printf("오류 발생 From db.Query: %v", err)
	}

	// email이 db에 없을 경우 db에 추가
	if !rows.Next() {
		db.Query("INSERT INTO user (email, provider) VALUES (?, ?)", payload.Email, loginRequest.Provider)
		pkg.BaseResponse(c, http.StatusOK, "success", payload)
		return
	}

	// email이 db에 있을 경우 accessToken, refreshToken 반환

	pkg.BaseResponse(c, http.StatusOK, "success", payload)
}

func login(c *gin.Context, redis *redis.Client, db *sql.DB) {
	return func(c *gin.Context) {

	}
}

// ID 토큰에서 사용자 이메일을 추출하는 함수
func GetUserEmailFromIdToken(c *gin.Context, redis *redis.Client, idToken string) (*Claims, error) {
	ISSUER := os.Getenv("KAKAO_ISSUER")
	KAKAO_REST_API_KEY := os.Getenv("KAKAO_REST_API_KEY")

	keys, err := GetPublicKeys(c, PROVIDER, redis)
	if err != nil {
		log.Printf("오류 발생 From GetPublicKeys: %v", err)
	}

	// idToken을 파싱하여 Header, Payload, Signature로 나누는 로직
	kid, err := getKidFromToken(idToken)
	if err != nil {
		log.Printf("오류 발생 From getKidFromToken: %v", err)
	}

	for _, key := range keys {
		if kid == key.Kid {
			// idToken을 파싱하여 Payload 추출
			publicKey, err := getRSAPublicKey(key)
			if err != nil {
				log.Printf("오류 발생 From getPayload: %v", err)
			}

			payload, err := validateSignature(idToken, publicKey, ISSUER, KAKAO_REST_API_KEY)
			if err != nil {
				log.Printf("오류 발생 From validateSignature: %v", err)
			}

			return payload, nil
		}
	}
	return nil, errors.New("유효한 키를 찾을 수 없음")
}

// JWT 토큰에서 kid 값을 추출하는 함수
func getKidFromToken(idToken string) (string, error) {
	header, err := getHeader(idToken)
	if err != nil {
		return "getHeader", err
	}

	decodedHeader, err := base64.RawURLEncoding.DecodeString(header)
	if err != nil {
		return "Base64", errors.New("Base64 디코딩 오류")
	}

	var headerJSON map[string]interface{}
	if err := json.Unmarshal(decodedHeader, &headerJSON); err != nil {
		return "JSON Parsing", errors.New("JSON 파싱 오류")
	}

	kid, ok := headerJSON["kid"].(string)
	if !ok {
		return "kid", errors.New("kid 값을 찾을 수 없음")
	}

	return kid, nil
}

func getHeader(idToken string) (string, error) {
	dividedToken, err := splitToken(idToken)
	if err != nil {
		return "", err
	}
	return dividedToken[0], nil
}

func splitToken(idToken string) ([]string, error) {
	dividedToken := strings.Split(idToken, ".")
	if len(dividedToken) != 3 {
		return nil, errors.New("JWT 토큰이 유효하지 않음")
	}
	return dividedToken, nil
}

// RSA 공개 키 생성 함수
func getRSAPublicKey(selectedKey JsonWebKey) (*rsa.PublicKey, error) {
	if selectedKey.Kty != "RSA" {
		return nil, errors.New("지원되지 않는 키 타입")
	}

	decodeM, err := base64UrlDecode(selectedKey.N)
	if err != nil {
		return nil, err
	}

	decodeE, err := base64UrlDecode(selectedKey.E)
	if err != nil {
		return nil, err
	}

	m := new(big.Int).SetBytes(decodeM)
	e := new(big.Int).SetBytes(decodeE).Int64()

	rsaPublicKey := &rsa.PublicKey{
		N: m,
		E: int(e),
	}

	return rsaPublicKey, nil
}

// JWT 서명 검증 함수
func validateSignature(idToken string, signingKey *rsa.PublicKey, issuer, audience string) (*Claims, error) {
	// 파서 설정
	token, err := jwt.ParseWithClaims(idToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("지원되지 않는 서명 방식")
		}
		return signingKey, nil
	})

	if err != nil {
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) {
			switch validationError.Errors {
			case jwt.ValidationErrorMalformed:
				return nil, errors.New("JWT가 올바르지 않음")
			case jwt.ValidationErrorSignatureInvalid:
				return nil, errors.New("서명이 유효하지 않음")
			case jwt.ValidationErrorExpired:
				return nil, errors.New("JWT가 만료됨")
			case jwt.ValidationErrorClaimsInvalid:
				return nil, errors.New("클레임이 유효하지 않음")
			default:
				return nil, err
			}
		}
		return nil, err
	}

	// 토큰 클레임 확인
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Issuer != issuer {
			return nil, fmt.Errorf("예상치 못한 발급자: %s", claims.Issuer)
		}
		if claims.Audience != audience {
			return nil, fmt.Errorf("예상치 못한 대상자: %s", claims.Audience)
		}
		return claims, nil
	}

	return nil, errors.New("클레임이 유효하지 않음")
}

// Base64 URL 디코딩 함수
func base64UrlDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
}

// 공개키 목록 조회 URL 요청 함수
func SetPublicKeys(c *gin.Context, redis *redis.Client) {
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
	key := redis.Set(c, PROVIDER, jsonData, 0)
	log.Println("데이터가 성공적으로 Redis에 저장되었습니다." + key.Val())

	pkg.BaseResponse(c, http.StatusOK, "공개키 저장 성공", key)
}

// Redis에서 공개키 가져오기 함수
func GetPublicKeys(c *gin.Context, provider string, redis *redis.Client) ([]JsonWebKey, error) {
	response := redis.Get(c, provider)
	if err := response.Err(); err != nil {
		log.Printf("오류 발생: %v", err)

		// 공개키 설정 함수 호출
		SetPublicKeys(c, redis)

		// 다시 시도하여 공개키 가져오기
		response = redis.Get(c, provider)
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
