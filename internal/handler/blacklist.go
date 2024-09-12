package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"
	"net/http"
)

type blacklistRequest struct {
	MemberId int64 `json:"memberId"`
}

// AddBlacklist godoc
// @Summary      memberId를 통해 차단하기
// @Description  memberId를 통해 차단하기
// @Tags         Blacklist
// @Accept       json
// @Produce      json
// @Param        blacklistRequest   body      blacklistRequest  true  "blacklistRequest"
// @Success      200 "성공"
// @Router       /v1/blacklist [post]
// @Security BearerAuth
func AddBlacklist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &blacklistRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		if request.MemberId == 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId is 0", nil)
			return
		}

		// 이미 차단했는지 확인
		isBlocked, err := mysql.Blacklists(qm.Where("blocker_member_id = ? AND blocked_member_id = ?", blockerId, request.MemberId)).Exists(c.Request.Context(), db)
		if isBlocked {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - already blocked", nil)
			return
		}
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		m := mysql.Blacklist{BlockerMemberID: blockerId.(int64), BlockedMemberID: request.MemberId}
		err = m.Insert(c.Request.Context(), db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

type deleteBlacklistRequest struct {
	MemberIds []int64 `json:"memberIds"`
}

// DeleteBlacklist godoc
// @Summary      memberId를 통해 차단해제하기
// @Description  memberId를 통해 차단해제하기
// @Tags         Blacklist
// @Accept       json
// @Produce      json
// @Param        deleteBlacklistRequest   body      deleteBlacklistRequest  true  "deleteBlacklistRequest"
// @Success      200 "성공"
// @Router       /v1/blacklist [delete]
// @Security BearerAuth
func DeleteBlacklist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &deleteBlacklistRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - "+err.Error(), nil)
			return
		}

		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		// Convert []int64 to []interface{}
		memberIDsInterface := make([]interface{}, len(request.MemberIds))
		for i, id := range request.MemberIds {
			memberIDsInterface[i] = id
		}

		log.Printf("memberIDsInterface: %v", memberIDsInterface)

		_, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId), qm.WhereIn("blocked_member_id IN ?", memberIDsInterface...)).DeleteAll(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		pkg.BaseResponse(c, http.StatusOK, "success", nil)
	}
}

type blacklistResponse struct {
	MemberId  int64  `json:"memberId"`
	Nickname  string `json:"nickname"`
	BlockDate string `json:"blockDate"`
}

// GetBlacklist godoc
// @Summary      memberId를 통해 차단 목록 조회하기
// @Description  memberId를 통해 차단 목록 조회하
// @Tags         Blacklist
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]blacklistResponse} "성공"
// @Router       /v1/blacklist [get]
// @Security BearerAuth
func GetBlacklist(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		blockerId, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - memberId not found", nil)
			return
		}

		all, err := mysql.Blacklists(qm.Where("blocker_member_id = ?", blockerId)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if len(all) == 0 {
			pkg.BaseResponse(c, http.StatusOK, "success", []blacklistResponse{})
			return
		}

		blockedIds := make([]interface{}, len(all))
		for i, entry := range all {
			blockedIds[i] = entry.BlockedMemberID
		}

		blockedMembers, err := mysql.Members(qm.WhereIn("member_id in ?", blockedIds...)).All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		blockedDateMap := make(map[int64]null.Time)
		for _, entry := range all {
			blockedDateMap[entry.BlockedMemberID] = entry.CreatedAt
		}

		responses := make([]blacklistResponse, 0, len(blockedMembers))
		for _, member := range blockedMembers {
			responses = append(responses, blacklistResponse{
				MemberId:  member.MemberID,
				Nickname:  member.Nickname.String,
				BlockDate: blockedDateMap[member.MemberID].Time.Format("2006-01-02 15:04:05"),
			})
		}

		pkg.BaseResponse(c, http.StatusOK, "success", responses)
	}
}
