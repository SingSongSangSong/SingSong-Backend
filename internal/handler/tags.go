package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tagMapToEnglish = map[string]string{
	"댄스본능":    "Dance Actively",
	"쉬어가는노래":  "I want to Rest",
	"헤어졌을때":   "when i Break up with someone who i care a lot",
	"신남과잔잔사이": "Between Exciting and Calm mood",
	"비":       "in rainday, i am sad",
	"크리스마스":   "Christmas for both children and adult",
	"저음노래":    "Bass songs that are easy",
	"고음노래":    "really hard Soprano",
	"썸":       "early stage of crush",
}

var tagMapToKorean = make(map[string]string)

// tagMapToEnglish를 기반으로 tagMapToKorean을 초기화. 프로그램 시작 시 자동실행
func init() {
	for korean, english := range tagMapToEnglish {
		tagMapToKorean[english] = korean
	}
}

// ListSsssTags godoc
// @Summary      ssss 태그 목록 가져오기
// @Description  ssss 태그 목록을 조회합니다.
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Success      200 {object} []string "성공"
// @Router       /tags/ssss [get]
func (handler *Handler) ListSsssTags(c *gin.Context) {
	tags := make([]string, 0, len(tagMapToEnglish))
	for tag := range tagMapToEnglish {
		tags = append(tags, tag)
	}
	c.JSON(http.StatusOK, BaseResponse{"ok", tags})
}

func mapTagKoreanToEnglish(koreanTag string) (string, error) {
	if englishTag, exists := tagMapToEnglish[koreanTag]; exists {
		return englishTag, nil
	}
	return "", errors.New("tag not found, tag cannot convert to english")
}

func mapTagEnglishToKorean(englishTag string) (string, error) {
	if koreanTag, exists := tagMapToKorean[englishTag]; exists {
		return koreanTag, nil
	}
	return "", errors.New("tag not found, tag cannot convert to english")
}
