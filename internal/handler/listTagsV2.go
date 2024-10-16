package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/friendsofgo/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tagToColumn = map[string]string{
	"그시절띵곡": "classics",
	"마무리송":  "finale",
	"고음송":   "high",
	"저음송":   "low",
	"댄스송":   "dance",
	"발라드송":  "ballads",
	"이별송":   "breakup",
	"R&B송":  "rnb",
	"듀엣송":   "duet",
	"썸송":    "ssum",
	"캐롤송":   "carol",
	"비올때송":  "rainy",
	"팝스타송":  "pop",
	"사회생활송": "office",
	"축가송":   "wedding",
	"입대송":   "military",
}

var defaultOrderV2 = []string{
	"그시절띵곡",
	"마무리송",
	"고음송",
	"저음송",
	"댄스송",
	"발라드송",
	"이별송",
	"R&B송",
	"비올때송",
	"캐롤송",
	"썸송",
	"듀엣송",
	"축가송",
	"입대송",
	"사회생활송",
	"팝스타송",
}

// ListTagsV2 godoc
// @Summary      ssss 태그 목록 가져오기 V2
// @Description  ssss 태그 목록을 조회합니다 V2
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct(data=[]string) "성공"
// @Router       /v2/tags [get]
func ListTagsV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		tags := make([]string, 0, len(defaultOrderV2))
		for _, tag := range defaultOrderV2 {
			tags = append(tags, tag)
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", tags)
	}
}

func MapTagToColumn(koreanTag string) (string, error) {
	if column, exists := tagToColumn[koreanTag]; exists {
		return column, nil
	}
	return "", errors.New("tag not found, tag cannot convert to database column:" + koreanTag)
}
