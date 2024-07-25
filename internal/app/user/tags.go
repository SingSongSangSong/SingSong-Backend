package user

import (
	"SingSong-Backend/internal/pkg"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var tagMapToEnglish = map[string]string{
	"Dance The Night Away": "Dance Actively",
	"잠깐! 쉬어가실게요~":          "I want to Rest",
	"이별후... 나는 가끔 눈물을 흘린다": "when i Break up with someone who i care a lot",
	"신남에 잔잔 두스푼":           "Between Exciting and Calm mood",
	"비도오고 그래서..":           "rainday i am sad",
	"산타도 인정한 캐롤송":          "Christmas for both children and adult",
	"음치 탈출 넘버원!":           "Bass songs that are easy",
	"지붕 뚫는 고음":             "really hard Soprano",
	"내꺼인듯 내꺼아닌 너":          "early stage of crush",
	"그 시절 띵곡":              "reminiscence of the past",
	"두근두근 듀엣송":             "Duet",
	"결혼 축하축가송~":            "celebrate the wedding",
	"필승! 입대를 명 받았습니다!":     "army",
	"마무리 1분 노래!":           "ending songs",
	"사회생활 S.O.S":           "for old people or senior",
	"내가바로팝스타":              "PopStar", //pop songs which are famous
}

var tagMapToKorean = make(map[string]string)

// tagMapToEnglish를 기반으로 tagMapToKorean을 초기화. 프로그램 시작 시 자동실행
func init() {
	for korean, english := range tagMapToEnglish {
		tagMapToKorean[english] = korean
	}

	log.Printf("init tagMapToKorean success")
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
	pkg.BaseResponse(c, http.StatusOK, "ok", tags)
	return
}

func MapTagKoreanToEnglish(koreanTag string) (string, error) {
	if englishTag, exists := tagMapToEnglish[koreanTag]; exists {
		return englishTag, nil
	}
	return "", errors.New("tag not found, tag cannot convert to english:" + koreanTag)
}

func MapTagEnglishToKorean(englishTag string) (string, error) {
	if koreanTag, exists := tagMapToKorean[englishTag]; exists {
		return koreanTag, nil
	}
	return "", errors.New("tag not found, tag cannot convert to korean:" + englishTag)
}

func MapTagsKoreanToEnglish(koreanTags []string) ([]string, error) {
	englishTags := make([]string, len(koreanTags))
	for i, tag := range koreanTags {
		englishTag, err := MapTagKoreanToEnglish(tag)
		if err != nil {
			return nil, err
		}
		englishTags[i] = englishTag
	}
	return englishTags, nil
}

func MapTagsEnglishToKorean(englishTags []string) ([]string, error) {
	koreanTags := make([]string, len(englishTags))
	for i, tag := range englishTags {
		koreanTag, err := MapTagEnglishToKorean(tag)
		if err != nil {
			return nil, err
		}
		koreanTags[i] = koreanTag
	}
	return koreanTags, nil
}
