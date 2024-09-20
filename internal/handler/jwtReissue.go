package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type ReissueRequest struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Reissue godoc
// @Summary      AccessToken RefreshToken 재발급
// @Description  AccessToken 재발급 및 RefreshToken 재발급 (RTR Refresh Token Rotation)
// @Tags         Reissue
// @Accept       json
// @Produce      json
// @Param        songs   body      ReissueRequest  true  "accessToken 및 refreshToken"
// @Success      200 {object} pkg.BaseResponseStruct{data=LoginResponse} "성공"
// @Router       /v1/member/reissue [post]
func Reissue(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		reissueRequest := &ReissueRequest{}
		if err := c.ShouldBindJSON(&reissueRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "JSON BINDING error - "+err.Error(), nil)
			return
		}
		// refreshToken이 redis에 있는지 확인
		_, err := redis.Get(c, reissueRequest.RefreshToken).Result()
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "No Refresh token in server - "+err.Error(), nil)
			return
		}

		// 만료된 accessToken에서 memberId를 추출 (서명 검증을 건너뛰고 페이로드만 파싱)
		token, _, err := new(jwt.Parser).ParseUnverified(reissueRequest.AccessToken, &Claims{})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "Token parsing error - "+err.Error(), nil)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "Invalid token claims", nil)
			return
		}

		// accessToken, refreshToken 생성
		accessTokenString, refreshTokenString, tokenErr := createAccessTokenAndRefreshToken(c, redis, claims, claims.BirthYear, claims.Gender, claims.MemberId)

		if tokenErr != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot create token "+tokenErr.Error(), nil)
			return
		}

		// 기존 refreshToken삭제
		_, err = redis.Del(c, reissueRequest.RefreshToken).Result()
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "Delete Redis error - "+err.Error(), nil)
			return
		}

		// JSON 응답 생성
		loginResponse := LoginResponse{
			AccessToken:  accessTokenString,
			RefreshToken: refreshTokenString,
		}

		// accessToken, refreshToken 반환
		pkg.BaseResponse(c, http.StatusOK, "success", loginResponse)
	}
}
