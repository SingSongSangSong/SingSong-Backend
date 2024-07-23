package user

import (
	"SingSong-Backend/config"
	"SingSong-Backend/internal/pkg"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
)

const (
	REQUEST_URL = "https://kauth.kakao.com/.well-known/jwks.json" // 공개키 목록 조회 URL
	PROVIDER    = "KAKAO"                                         // 공급자
)

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
	IdToken string `json:"IdToken"`
}

// ID 토큰에서 사용자 이메일을 추출하는 함수
func (handler *RedisHandler) GetUserEmailFromIdToken(c *gin.Context) {
	loginRequest := &LoginRequest{}
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
		return
	}

	jwtConfig := config.NewJwtConfig(os.Getenv("KAKAO_ISSUER"), os.Getenv("KAKAO_REST_API_KEY"))
	ISSUER := jwtConfig.Issuer
	KAKAO_REST_API_KEY := jwtConfig.Audience

	keys, err := handler.GetPublicKeys(c, PROVIDER)
	if err != nil {
		log.Printf("오류 발생 From GetPublicKeys: %v", err)
	}

	// idToken을 파싱하여 Header, Payload, Signature로 나누는 로직
	kid, err := getKidFromToken(loginRequest.IdToken)
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

			payload, err := validateSignature(loginRequest.IdToken, publicKey, ISSUER, KAKAO_REST_API_KEY)
			if err != nil {
				log.Printf("오류 발생 From validateSignature: %v", err)
			}

			pkg.BaseResponse(c, http.StatusOK, "사용자 이메일 추출 성공", payload.Email)
		}
	}
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
