package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"strconv"
)

type LoginV2Request struct {
	IdToken  string `json:"idToken"`
	Provider string `json:"provider"`
}

type LoginV2Response struct {
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	IsInfoRequired bool   `json:"isInfoRequired"`
}

// LoginV2 로그인 API
// @Summary 로그인 API
// @Description 로그인 API
// @Tags Signup and Login
// @Accept json
// @Produce json
// @Param loginV2 body LoginV2Request true "로그인 요청"
// @Success 200 {object} pkg.BaseResponseStruct{data=LoginV2Response} "로그인 성공"
// @Router  /v2/member/login [post]
func LoginV2(rdb *redis.Client, db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		loginRequest := &LoginV2Request{}
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

			accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, rdb, &Claims{Email: "Anonymous@anonymous.com"}, "0", "Unknown", m.MemberID)
			if tokenErr != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
				return
			}

			loginResponse := LoginV2Response{
				AccessToken:    accessTokenString,
				RefreshToken:   refreshTokenString,
				IsInfoRequired: false,
			}

			// accessToken, refreshToken 반환
			pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
			return
		}

		// email 및 nickname 추출
		payload, err := GetUserEmailFromIdToken(c, rdb, loginRequest.IdToken, loginRequest.Provider)
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
			m, err = joinV2(c, payload, loginRequest, m, db)
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
			go CreatePlaylist(db, m.Nickname.String+null.StringFrom("의 플레이리스트").String, m.MemberID)

			accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, rdb, payload, strconv.Itoa(m.Birthyear.Int), m.Gender.String, m.MemberID)

			if tokenErr != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
				return
			}

			loginResponse := LoginV2Response{
				AccessToken:    accessTokenString,
				RefreshToken:   refreshTokenString,
				IsInfoRequired: true,
			}

			// accessToken, refreshToken 반환
			pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
			return
		}

		accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, rdb, payload, strconv.Itoa(m.Birthyear.Int), m.Gender.String, m.MemberID)

		if tokenErr != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
			return
		}

		loginResponse := LoginV2Response{
			AccessToken:    accessTokenString,
			RefreshToken:   refreshTokenString,
			IsInfoRequired: false,
		}

		// accessToken, refreshToken 반환
		pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
	}
}

type LoginV2ExtraInfoRequest struct {
	BirthYear string `json:"birthYear"`
	Gender    string `json:"gender"`
}

// LoginV2ExtraInfoRequired
// @Summary 로그인 성별 및 연령 정보가 필요할때 사용, InfoRequired가 true일때만 사용
// @Description 로그인 성별 및 연령 정보 받는 API
// @Tags Signup and Login
// @Accept json
// @Produce json
// @Param loginV2 body LoginV2ExtraInfoRequest true "로그인 요청"
// @Success 200 {object} pkg.BaseResponseStruct{data=nil}  "로그인 성공"
// @Router  /v2/member/login/extra [post]
// @Security BearerAuth
func LoginV2ExtraInfoRequired(db *sql.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		loginExtraInfoRequest := &LoginV2ExtraInfoRequest{}
		if err := c.ShouldBindJSON(&loginExtraInfoRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		birthYearInt, err := strconv.Atoi(loginExtraInfoRequest.BirthYear)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// soft delete
		_, err = mysql.Members(
			qm.Where("member_id = ?", memberId), qm.And("deleted_at IS NULL"),
		).UpdateAll(c.Request.Context(), db, mysql.M{"birthyear": null.IntFrom(birthYearInt), "gender": null.StringFrom(loginExtraInfoRequest.Gender)})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

func joinV2(c *gin.Context, payload *Claims, loginRequest *LoginV2Request, m *mysql.Member, db *sql.DB) (*mysql.Member, error) {
	// 랜덤 닉네임
	nickname := generateRandomNickname()
	nullNickname := null.StringFrom(nickname)
	m = &mysql.Member{Provider: loginRequest.Provider, Email: payload.Email, Nickname: nullNickname}
	err := m.Insert(c.Request.Context(), db, boil.Infer())
	if err != nil {
		return nil, errors.New("error inserting member - " + err.Error())
	}
	return m, nil
}
