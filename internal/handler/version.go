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
	UpdateUrl          string `json:"updateUrl"`
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
			UpdateUrl:          version.UpdateURL,
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

type VersionUpdateRequest struct {
	LatestVersion      string `json:"latestVersion"`
	ForceUpdateVersion string `json:"forceUpdateVersion,omitempty"`
	Platform           string `json:"platform"`
	UpdateUrl          string `json:"updateUrl,omitempty"`
}

// VersionUpdate godoc
// @Summary      버전 추가
// @Description  플랫폼별 최신 버전, 강제 업데이트 버전을 설정할 수 있다. (강제 업데이트 버전이랑 update url은 걍 필드 빼고 보내면 갱신 안됨)
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
		one, err := mysql.AppVersions(qm.Where("platform = ?", request.Platform)).One(c.Request.Context(), db)
		if err != nil {
			// row가 없는 경우 새로운 레코드 추가
			if err == sql.ErrNoRows {
				// 강제 업데이트 버전과 업데이트 URL이 필요한 경우 확인
				if request.ForceUpdateVersion == "" || request.UpdateUrl == "" {
					pkg.BaseResponse(c, http.StatusBadRequest, "error - force update version and update URL are required", nil)
					return
				}
				// 새 버전 정보 추가
				version := mysql.AppVersion{
					Platform:           request.Platform,
					LatestVersion:      request.LatestVersion,
					ForceUpdateVersion: request.ForceUpdateVersion,
					UpdateURL:          request.UpdateUrl,
				}
				if err := version.Insert(c.Request.Context(), db, boil.Infer()); err != nil {
					pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
					return
				}
			} else {
				// 그 외 데이터베이스 오류 처리
				pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
				return
			}
		} else {
			// 기존 row가 있는 경우 업데이트
			one.LatestVersion = request.LatestVersion
			if request.ForceUpdateVersion != "" {
				one.ForceUpdateVersion = request.ForceUpdateVersion
			}
			if request.UpdateUrl != "" {
				one.UpdateURL = request.UpdateUrl
			}
			// 업데이트 실행
			if _, err := one.Update(c.Request.Context(), db, boil.Infer()); err != nil {
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
				UpdateUrl:          v.UpdateURL,
			})
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}
