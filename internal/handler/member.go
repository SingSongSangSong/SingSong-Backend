package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"time"
)

type MemberResponse struct {
	Nickname  string `json:"nickname"`
	Birthyear int    `json:"birthYear"`
	Gender    string `json:"gender"`
	Email     string `json:"email"`
}

// GetMemberInfo godoc
// @Summary      Member의 정보를 가져온다
// @Description  사용자 정보 조회
// @Tags         Member
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=MemberResponse} "성공"
// @Router       /member [get]
// @Security BearerAuth
func GetMemberInfo(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get memberId from context
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Get member info
		member, err := mysql.Members(qm.Where("member_id = ?", memberId)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		nullNickname := member.Nickname.String
		nullBirthyear := member.Birthyear.Int
		nullGender := member.Gender.String

		// JSON response
		memberResponse := MemberResponse{
			Email:     member.Email,
			Nickname:  nullNickname,
			Birthyear: nullBirthyear,
			Gender:    nullGender,
		}

		pkg.BaseResponse(c, http.StatusOK, "success", memberResponse)
	}
}

type WithdrawRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// Withdraw godoc
// @Summary      멤버 회원 탈퇴
// @Description  멤버 회원 탈퇴
// @Tags         Member
// @Accept       json
// @Produce      json
// @Param        refreshToken   body      WithdrawRequest  true  "refreshToken"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /member/withdraw [post]
// @Security BearerAuth
func Withdraw(db *sql.DB, redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		withdrawRequest := &WithdrawRequest{}
		if err := c.ShouldBindJSON(&withdrawRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// Get memberId from context
		memberId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Delete member
		_, err := mysql.Members(qm.Where("member_id = ?", memberId)).UpdateAll(c, db, mysql.M{"deleted_at": time.Now()})
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		// Delete redis
		_, err = redis.Del(c, withdrawRequest.RefreshToken).Result()
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

// Logout godoc
// @Summary      멤버 회원 로그아웃
// @Description  멤버 회원 로그아웃
// @Tags         Member
// @Accept       json
// @Produce      json
// @Param        refreshToken   body      WithdrawRequest  true  "refreshToken"
// @Success      200 {object} pkg.BaseResponseStruct{} "성공"
// @Router       /member/logout [post]
// @Security BearerAuth
func Logout(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		withdrawRequest := &WithdrawRequest{}
		if err := c.ShouldBindJSON(&withdrawRequest); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		// Delete redis
		_, err := redis.Del(c, withdrawRequest.RefreshToken).Result()
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}
