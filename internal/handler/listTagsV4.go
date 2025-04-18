package handler

import (
	"SingSong-Server/internal/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

// v2와 똑같은데 태그 순서만 바뀜
var tagColumnMappingsV4 = []TagColumnMapping{
	{Tag: "밴드송", Column: "band"},
	{Tag: "힙합송", Column: "hiphop"},
	{Tag: "J-POP송", Column: "jpop"},
	{Tag: "뮤지컬송", Column: "musical"},
	{Tag: "그시절띵곡", Column: "classics"},
	{Tag: "마무리송", Column: "finale"},
	{Tag: "고음송", Column: "high"},
	{Tag: "저음송", Column: "low"},
	{Tag: "댄스송", Column: "dance"},
	{Tag: "발라드송", Column: "ballads"},
	{Tag: "이별송", Column: "breakup"},
	{Tag: "R&B송", Column: "rnb"},
	{Tag: "듀엣송", Column: "duet"},
	{Tag: "썸송", Column: "ssum"},
	{Tag: "팝스타송", Column: "pop"},
	{Tag: "비올때송", Column: "rainy"},
	{Tag: "캐롤송", Column: "carol"},
	{Tag: "사회생활송", Column: "office"},
	{Tag: "축가송", Column: "wedding"},
	{Tag: "입대송", Column: "military"},
}

// 태그에서 컬럼으로 빠르게 접근할 수 있도록 맵을 생성
var tagToColumnV4 = make(map[string]string)

func init() {
	// tagToColumn 맵을 초기화 (한 번만 실행)
	for _, mapping := range tagColumnMappingsV4 {
		tagToColumnV4[mapping.Tag] = mapping.Column
	}
}

// ListTagsV4 godoc
// @Summary      태그 목록 가져오기 V4 (v2와 팝스타송, 캐롤송 순서가 바뀜 + 뮤지컬/밴드/jpop/힙합 추가)
// @Description  태그 목록을 조회합니다 V4
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct(data=[]string) "성공"
// @Router       /v4/tags [get]
func ListTagsV4() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 태그 목록 반환 (기본 정렬된 순서로)
		tags := make([]string, 0, len(tagColumnMappingsV4))
		for _, mapping := range tagColumnMappingsV4 {
			tags = append(tags, mapping.Tag)
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", tags)
	}
}

func MapTagToColumnV4(koreanTag string) (string, error) {
	if column, exists := tagToColumnV4[koreanTag]; exists {
		return column, nil
	}
	return "", errors.Wrap(fmt.Errorf("tag not found, tag cannot convert to english: %s", koreanTag), "최초 에러 발생 지점")
}
