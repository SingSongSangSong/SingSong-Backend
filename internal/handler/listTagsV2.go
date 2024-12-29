package handler

import (
	"SingSong-Server/internal/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type TagColumnMapping struct {
	Tag    string // 태그 이름 (한글)
	Column string // 데이터베이스 컬럼 이름
}

// 태그와 컬럼을 하나의 구조체로 관리
var tagColumnMappings = []TagColumnMapping{
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
	{Tag: "캐롤송", Column: "carol"},
	{Tag: "비올때송", Column: "rainy"},
	{Tag: "팝스타송", Column: "pop"},
	{Tag: "사회생활송", Column: "office"},
	{Tag: "축가송", Column: "wedding"},
	{Tag: "입대송", Column: "military"},
}

// 태그에서 컬럼으로 빠르게 접근할 수 있도록 맵을 생성
var tagToColumn = make(map[string]string)

func init() {
	// tagToColumn 맵을 초기화 (한 번만 실행)
	for _, mapping := range tagColumnMappings {
		tagToColumn[mapping.Tag] = mapping.Column
	}
}

// ListTagsV2 godoc
// @Summary      태그 목록 가져오기 V2
// @Description  태그 목록을 조회합니다 V2
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct(data=[]string) "성공"
// @Router       /v2/tags [get]
func ListTagsV2() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 태그 목록 반환 (기본 정렬된 순서로)
		tags := make([]string, 0, len(tagColumnMappings))
		for _, mapping := range tagColumnMappings {
			tags = append(tags, mapping.Tag)
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", tags)
	}
}

// MapTagToColumn - 태그를 데이터베이스 컬럼으로 매핑
func MapTagToColumn(koreanTag string) (string, error) {
	if column, exists := tagToColumn[koreanTag]; exists {
		return column, nil
	}
	return "", errors.Wrap(fmt.Errorf("tag not found, tag cannot convert to database column: "+koreanTag), "최초 에러 발생 지점")
}
