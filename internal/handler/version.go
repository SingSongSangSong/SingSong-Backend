package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"net/http"
)

type versionCheckRequest struct {
	Platform string `json:"platform"`
}

type versionCheckResponse struct {
	LatestVersion      string `json:"latestVersion"`
	ForceUpdateVersion string `json:"forceUpdateVersion"`
	Platform           string `json:"platform"`
}

// VersionCheck godoc
// @Summary      버전 확인
// @Description  request body에 플랫폼을 보내면, 해당 플랫폼의 최신 버전과 강제 업데이트 버전을 응답
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Param        version  body      versionCheckRequest  true  "플랫폼 정보"
// @Success      200 {object} pkg.BaseResponseStruct(data=versionCheckResponse) "성공"
// @Router       /v1/version/check [post]
func VersionCheck(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &versionCheckRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no version in request", nil)
			return
		}

		version, err := mysql.AppVersions(qm.Where("platform = ?", request.Platform)).One(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - cannot find version data", nil)
			return
		}

		response := &versionCheckResponse{
			LatestVersion:      version.LatestVersion,
			ForceUpdateVersion: version.ForceUpdateVersion,
			Platform:           version.Platform,
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type VersionUpdateRequest struct {
	LatestVersion      string `json:"latestVersion"`
	ForceUpdateVersion string `json:"forceUpdateVersion"`
	Platform           string `json:"platform"`
}

// VersionUpdate godoc
// @Summary      버전 추가
// @Description  플랫폼별 최신 버전, 강제 업데이트 버전을 설정할 수 있다. (강제 업데이트 버전을 빈 문자열로 보내면 강제 업데이트 버전은 갱신안됨)
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Param 		version body      VersionUpdateRequest  true  "등록 버전 정보"
// @Success      200 {object} pkg.BaseResponseStruct(data=nil) "성공"
// @Router       /v1/version/update [post]
func VersionUpdate(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &VersionUpdateRequest{}
		if err := c.ShouldBindJSON(&request); err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - no version in request", nil)
			return
		}
		// db에 정보를 조회
		exists, err := mysql.AppVersions(qm.Where("platform = ?", request.Platform)).Exists(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error -"+err.Error(), nil)
			return
		}
		if exists {
			// 기존 버전 정보가 존재할 경우 업데이트
			// 강제 업데이트 버전이 빈 문자열이 아닌 경우에만 업데이트
			if request.ForceUpdateVersion != "" {
				_, err := mysql.AppVersions(
					qm.Where("platform = ?", request.Platform),
				).UpdateAll(c, db, mysql.M{
					"latest_version":       request.LatestVersion,
					"force_update_version": request.ForceUpdateVersion,
				})
				if err != nil {
					pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
					return
				}
			} else {
				// 강제 업데이트 버전이 빈 문자열일 경우, 최신 버전만 업데이트
				_, err := mysql.AppVersions(
					qm.Where("platform = ?", request.Platform),
				).UpdateAll(c, db, mysql.M{
					"latest_version": request.LatestVersion,
				})
				if err != nil {
					pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
					return
				}
			}
		} else {
			//새로 추가
			var forceUpdateVersion = request.ForceUpdateVersion
			if forceUpdateVersion == "" {
				forceUpdateVersion = "1.0.0"
			}
			version := mysql.AppVersion{Platform: request.Platform, LatestVersion: request.LatestVersion, ForceUpdateVersion: forceUpdateVersion}
			err := version.Insert(c, db, boil.Infer())
			if err != nil {
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", nil)
	}
}

// AllVersion godoc
// @Summary      모든 버전 확인
// @Description  플랫폼별 등록되어 있는 버전 내용 확인 가능
// @Tags         App Version
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct{data=[]versionCheckResponse} "성공"
// @Router       /v1/version [get]
func AllVersion(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		all, err := mysql.AppVersions().All(c.Request.Context(), db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}

		var response []versionCheckResponse
		for _, v := range all {
			response = append(response, versionCheckResponse{
				Platform:           v.Platform,
				LatestVersion:      v.LatestVersion,
				ForceUpdateVersion: v.ForceUpdateVersion,
			})
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
