package handler

import (
	"SingSong-Server/conf"
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const (
	//kakao
	KAKAO_PUBLIC_KEY_URL = "https://kauth.kakao.com/.well-known/jwks.json"
	KAKAO_PROVIDER       = "KAKAO_KEY"
	//apple
	APPLE_PUBLIC_KEY_URL = "https://appleid.apple.com/auth/keys"
	APPLE_PROVIDER       = "APPLE_KEY"
)

var (
	SECRET_KEY                   = conf.AuthConfigInstance.SECRET_KEY
	KAKAO_REST_API_KEY           = conf.AuthConfigInstance.KAKAO_REST_API_KEY
	KAKAO_ISSUER                 = conf.AuthConfigInstance.KAKAO_ISSUER
	APPLE_CLIENT_ID              = conf.AuthConfigInstance.APPLE_CLIENT_ID
	APPLE_ISSUER                 = conf.AuthConfigInstance.APPLE_ISSUER
	JWT_ISSUER                   = conf.AuthConfigInstance.JWT_ISSUER
	JWT_ACCESS_VALIDITY_SECONDS  = conf.AuthConfigInstance.JWT_ACCESS_VALIDITY_SECONDS
	JWT_REFRESH_VALIDITY_SECONDS = conf.AuthConfigInstance.JWT_REFRESH_VALIDITY_SECONDS
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
	MemberId  int64  `json:"memberId"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Gender    string `json:"gender"`
	BirthYear string `json:"birthYear"`
	Picture   string `json:"picture"`
	Provider  string `json:"provider"`
	jwt.StandardClaims
}

// 각 Key 구조체를 담을 구조체 정의
type KeyContainer struct {
	Keys []JsonWebKey `json:"keys"`
}

type LoginRequest struct {
	IdToken   string `json:"idToken"`
	Provider  string `json:"provider"`
	BirthYear string `json:"birthYear"`
	Gender    string `json:"gender"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// 고유한 이메일 생성 함수
func generateUniqueEmail() string {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to generate UUID: %v", err)
	}
	// UUID를 기반으로 이메일 생성
	return fmt.Sprintf("Anonymous+%s@anonymous.com", newUUID.String())
}

// Login godoc
// @Summary      회원가입 및 로그인
// @Description  IdToken을 이용한 회원가입 및 로그인
// @Tags         Signup and Login
// @Accept       json
// @Produce      json
// @Param        songs   body      LoginRequest  true  "idToken 및 Provider"
// @Success      200 {object} pkg.BaseResponseStruct{data=LoginResponse} "성공"
// @Router       /v1/member/login [post]
func Login(redis *redis.Client, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		loginRequest := &LoginRequest{}
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// provider가 Anonymous인 경우
		if loginRequest.Provider == "Anonymous" {
			// DB에 없는 경우 - 회원가입
			m, err := joinForAnonymous(c, &Claims{Email: generateUniqueEmail()}, 0, "Unknown", loginRequest.Provider, db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			go CreatePlaylist(db, m.Nickname.String+null.StringFrom("의 플레이리스트").String, m.MemberID)

			accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, redis, &Claims{Email: "Anonymous@anonymous.com"}, "0", "Unknown", m.MemberID)
			if tokenErr != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
				return
			}

			loginResponse := LoginResponse{
				AccessToken:  accessTokenString,
				RefreshToken: refreshTokenString,
			}

			// accessToken, refreshToken 반환
			pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
			return
		}

		// email 및 nickname 추출
		payload, err := GetUserEmailFromIdToken(c, redis, loginRequest.IdToken, loginRequest.Provider)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// email이 없을 경우 에러 반환
		if payload.Email == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - Email is empty", nil)
			return
		}

		// email+Provider db에 있는지 확인
		m, err := mysql.Members(qm.Where("email = ? AND provider = ? AND deleted_at is null", payload.Email, loginRequest.Provider)).One(c.Request.Context(), db)
		if err != nil {
			// DB에 없는 경우 - 회원가입
			m, err = join(c, payload, loginRequest, m, db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			go CreatePlaylist(db, m.Nickname.String+null.StringFrom("의 플레이리스트").String, m.MemberID)
		}

		accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, redis, payload, strconv.Itoa(m.Birthyear.Int), m.Gender.String, m.MemberID)

		if tokenErr != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
			return
		}

		loginResponse := LoginResponse{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		}

		// accessToken, refreshToken 반환
		pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
	}
}

// GetUserEmailFromIdToken ID 토큰에서 사용자 이메일을 추출하는 함수
func GetUserEmailFromIdToken(c *gin.Context, redis *redis.Client, idToken string, provider string) (*Claims, error) {
	var issuer, apiKey string
	switch provider {
	case KAKAO_PROVIDER:
		issuer = KAKAO_ISSUER
		apiKey = KAKAO_REST_API_KEY
	case APPLE_PROVIDER:
		issuer = APPLE_ISSUER
		apiKey = APPLE_CLIENT_ID
	default:
		return nil, errors.New("유효하지 않은 provider")
	}

	keys, err := GetPublicKeys(c, provider, redis)
	if err != nil {
		log.Printf("오류 발생 From GetPublicKeys: %v", err)
		return nil, err
	}

	// idToken을 파싱하여 Header, Payload, Signature로 나누는 로직
	kid, err := getKidFromToken(idToken)
	if err != nil {
		log.Printf("오류 발생 From getKidFromToken: %v", err)
		return nil, err
	}

	for _, key := range keys {
		if kid == key.Kid {
			publicKey, err := getRSAPublicKey(key)
			if err != nil {
				log.Printf("오류 발생 From getPayload: %v", err)
				return nil, err
			}

			payload, err := validateSignature(idToken, publicKey, issuer, apiKey)
			if err != nil {
				log.Printf("오류 발생 From validateSignature: %v", err)
				return nil, err
			}

			return payload, nil
		}
	}
	return nil, errors.New("유효한 키를 찾을 수 없음")
}

// Redis에서 공개키 가져오기 함수
func GetPublicKeys(c *gin.Context, provider string, redis *redis.Client) ([]JsonWebKey, error) {
	response := redis.Get(c, provider)

	// redis에 공개키가 없는 경우에 public key url을 통해 가져와 저장한다
	if err := response.Err(); err != nil {
		log.Printf("오류 발생: %v", err)

		//provider에 따른 public key url 설정
		var publicKeyUrl string
		switch provider {
		case KAKAO_PROVIDER:
			publicKeyUrl = KAKAO_PUBLIC_KEY_URL
		case APPLE_PROVIDER:
			publicKeyUrl = APPLE_PUBLIC_KEY_URL
		}

		// 공개키 설정 함수 호출
		err = FetchPublicKeys(c, redis, publicKeyUrl, provider)
		if err != nil {
			log.Printf("GetKaKaoPublicKey 오류 발생: %v", err)
			return nil, err
		}

		// 다시 시도하여 공개키 가져오기
		response = redis.Get(c, provider)
		if err := response.Err(); err != nil {
			log.Printf("오류 발생 (다시 시도 후): %v", err)
			return nil, err
		}
	}

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

// 외부 api로부터 공개키 목록 조회해서 Redis에 저장
func FetchPublicKeys(c *gin.Context, redis *redis.Client, publicKeyUrl string, provider string) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(publicKeyUrl)
	if err != nil {
		log.Printf("HTTP 요청 오류: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		log.Printf("HTTP 응답 오류: 상태 코드 %d", resp.StatusCode)
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("응답 본문 읽기 오류: %v", err)
		return err
	}

	// PublicKeyDto 구조체 생성
	publicKey := PublicKeyDto{
		Provider: provider,
		Key:      string(body),
	}

	// JSON으로 변환
	jsonData, err := json.Marshal(publicKey)
	if err != nil {
		log.Printf("JSON 마샬 오류: %v", err)
		return err
	}

	// Redis에 저장
	key := redis.Set(c, provider, jsonData, 24*time.Hour)
	log.Println("데이터가 성공적으로 Redis에 저장되었습니다." + key.Val())

	return nil
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

// JWT 토큰에서 kid 값을 추출하는 함수
func getKidFromToken(idToken string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		log.Printf("Error parsing kid from id token: %v\n", err)
		return "", err
	}
	kid := token.Header["kid"].(string)
	return kid, nil
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

// Base64 URL 디코딩 함수
func base64UrlDecode(data string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(data)
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

func joinForAnonymous(c *gin.Context, payload *Claims, year int, gender string, provider string, db *sql.DB) (*mysql.Member, error) {
	m := &mysql.Member{Provider: provider, Email: payload.Email, Nickname: null.StringFrom("Anonymous"), Birthyear: null.IntFrom(year), Gender: null.StringFrom(gender)}
	err := m.Insert(c.Request.Context(), db, boil.Infer())
	if err != nil {
		return nil, errors.New("error inserting member - " + err.Error())
	}
	return m, nil
}

func join(c *gin.Context, payload *Claims, loginRequest *LoginRequest, m *mysql.Member, db *sql.DB) (*mysql.Member, error) {
	// 랜덤 닉네임
	nickname := generateRandomNickname()
	nullNickname := null.StringFrom(nickname)

	// Convert the BirthYear to an integer
	birthYearInt, err := strconv.Atoi(loginRequest.BirthYear)
	if err != nil {
		log.Printf("Invalid BirthYear format: %v", err)
		birthYearInt = 0 // Set to 0 or handle appropriately
	}

	// Initialize the null.Int from the converted integer
	nullBrithyear := null.IntFrom(birthYearInt)
	nullGender := null.StringFrom(loginRequest.Gender)

	m = &mysql.Member{Provider: loginRequest.Provider, Email: payload.Email, Nickname: nullNickname, Birthyear: nullBrithyear, Gender: nullGender}
	err = m.Insert(c.Request.Context(), db, boil.Infer())
	if err != nil {
		//pkg.BaseResponse(c, http.StatusBadRequest, "error inserting member - "+err.Error(), nil)
		return nil, errors.New("Error inserting member")
	}
	return m, nil
}

// 랜덤 닉네임 제조기
var (
	firstPart  = []string{"귀여운", "멋쟁이", "행복한", "슬픈", "도도한", "스윗한", "차가운"}
	secondPart = []string{"고양이", "강아지", "토끼", "여우", "곰", "사자", "호랑이", "부엉이", "펭귄", "코끼리"}
)

func generateRandomNickname() string {
	first := firstPart[rand.Intn(len(firstPart))]
	second := secondPart[rand.Intn(len(secondPart))]
	return first + " " + second
}

func createAccessTokenAndRefreshToken(c *gin.Context, redis *redis.Client, payload *Claims, birthYear string, gender string, memberId int64) (string, string, error) {
	jwtAccessValidityStr := JWT_ACCESS_VALIDITY_SECONDS
	if jwtAccessValidityStr == "" {
		log.Printf("JWT_ACCESS_VALIDITY_SECONDS 환경 변수가 설정되지 않았습니다.")
		return "", "", fmt.Errorf("JWT_ACCESS_VALIDITY_SECONDS 환경 변수가 설정되지 않았습니다")
	}

	jwtAccessValidity, err := strconv.ParseInt(jwtAccessValidityStr, 10, 64)
	if err != nil {
		log.Printf("환경 변수 변환 실패: %v", err)
		return "", "", fmt.Errorf("환경 변수 변환 실패: %v", err)
	}

	accessTokenExpiresAt := time.Now().Add(time.Duration(jwtAccessValidity) * time.Second).Unix()
	at := Claims{
		MemberId:  memberId,
		Gender:    gender,
		BirthYear: birthYear,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpiresAt,
			Issuer:    JWT_ISSUER,
			IssuedAt:  time.Now().Unix(),
			Subject:   "AccessToken",
		},
	}

	jwtRefreshValidityStr := JWT_REFRESH_VALIDITY_SECONDS
	if jwtRefreshValidityStr == "" {
		log.Printf("JWT_REFRESH_VALIDITY_SECONDS 환경 변수가 설정되지 않았습니다.")
		return "", "", errors.New("JWT_REFRESH_VALIDITY_SECONDS 환경 변수가 설정되지 않았습니다")
	}

	jwtRefreshValidity, err := strconv.ParseInt(jwtRefreshValidityStr, 10, 64)
	if err != nil {
		log.Printf("환경 변수 변환 실패: %v", err)
		return "", "", fmt.Errorf("환경 변수 변환 실패: %v", err)
	}

	refreshTokenExpiresAt := time.Now().Add(time.Duration(jwtRefreshValidity) * time.Second).Unix()
	rt := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiresAt,
			Issuer:    JWT_ISSUER,
			IssuedAt:  time.Now().Unix(),
			Subject:   "RefreshToken",
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, at)
	accessTokenString, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, rt)
	refreshTokenString, err := refreshToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	payload.BirthYear = birthYear
	payload.Gender = gender
	payload.MemberId = memberId

	claims, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	_, err = redis.Set(c, refreshTokenString, claims, time.Duration(jwtRefreshValidity)*time.Second).Result()
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
