package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
)

// v2와 똑같은데 태그 순서만 바뀜
var tagColumnMappingsV3 = []TagColumnMapping{
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

// ListTagsV2 godoc
// @Summary      태그 목록 가져오기 V3 (v2와 팝스타송, 캐롤송 순서가 바뀜)
// @Description  태그 목록을 조회합니다 V3
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Success      200 {object} pkg.BaseResponseStruct(data=[]string) "성공"
// @Router       /v3/tags [get]
func ListTagsV3() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 태그 목록 반환 (기본 정렬된 순서로)
		tags := make([]string, 0, len(tagColumnMappingsV3))
		for _, mapping := range tagColumnMappingsV3 {
			tags = append(tags, mapping.Tag)
		}
		pkg.BaseResponse(c, http.StatusOK, "ok", tags)
	}
}
