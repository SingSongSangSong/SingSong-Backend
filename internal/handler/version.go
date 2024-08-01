package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
	"time"
)

type versionCheckRequest struct {
	Version string `json:"version"`
}

type versionCheckResponse struct {
	IsLatest    bool `json:"isLatest"`
	ForceUpdate bool `json:"forceUpdate"`
}

// VersionCheck godoc
// @Summary      버전 확인
// @Description  헤더에 플랫폼 정보를 포함하고, request body 앱의 버전을 보내면, 최신 버전인지 여부와 강제 업데이트 필요 여부를 응답
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Param        version  body      versionCheckRequest  true  "현재 앱 버전 정보"
// @Success      200 {object} pkg.BaseResponseStruct(data=versionCheckResponse) "성공"
// @Router       /version/check [post]
func VersionCheck(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("platform")
		if !exists {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - platform not found", nil)
			return
		}

		platform := value.(string)
		if platform != "android" && platform != "ios" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - platform not supported", nil)
			return
		}

		request := &versionCheckRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no version in request", nil)
			return
		}

		requestVersion, err := mysql.AppVersions(qm.Where("platform = ? AND version = ?", platform, request.Version)).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot find version data", nil)
			return
		}

		latestVersion, err := mysql.AppVersions(qm.Where("platform = ?", platform), qm.OrderBy("releaseDate DESC")).One(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		response := &versionCheckResponse{
			IsLatest:    latestVersion.Version == request.Version,
			ForceUpdate: requestVersion.ForceUpdate,
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type latestVersionUpdateRequest struct {
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	ForceUpdate bool   `json:"forceUpdate"`
}

// LatestVersionUpdate godoc
// @Summary      버전 추가
// @Description  새로운 버전이 나왔을때 버전을 추가할 수 있음 (플랫폼(ios, android), 버전, 이전 버전들을 강제 업데이트 할지 여부)
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Param 		version body      latestVersionUpdateRequest  true  "등록 버전 정보"
// @Success      200 "성공"
// @Router       /version/update [post]
func LatestVersionUpdate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &latestVersionUpdateRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no version in request", nil)
			return
		}
		// db에 정보를 저장한다
		version := mysql.AppVersion{Platform: request.Platform, Version: request.Version, ForceUpdate: false}
		err := version.Insert(c, db, boil.Infer())
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		if request.ForceUpdate == true {
			// 강제 업데이트가 필요하다면, 이전 버전들을 모두 강제 업데이트로 변경한다
			_, err := mysql.AppVersions(qm.Where("platform = ? AND version != ?", request.Platform, request.Version)).UpdateAll(c, db, mysql.M{"forceUpdate": true})
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}

type VersionResponse struct {
	Platform    string `json:"platform"`
	Version     string `json:"version"`
	ForceUpdate bool   `json:"forceUpdate"`
	ReleaseDate string `json:"releaseDate"`
}

// AllVersion godoc
// @Summary      모든 버전 확인
// @Description  등록되어 있는 모든 버전 확인 가능
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Success      200 "성공" {object} pkg.BaseResponseStruct{data=[]versionResponse} "성공"
// @Router       /version [get]
func AllVersion(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		all, err := mysql.AppVersions().All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		var response []VersionResponse
		for _, v := range all {
			response = append(response, VersionResponse{
				Platform:    v.Platform,
				Version:     v.Version,
				ForceUpdate: v.ForceUpdate,
				ReleaseDate: v.CreatedAt.Time.Format(time.RFC3339),
			})
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
