package handler

import (
	"SingSong-Server/internal/db/mysql"
	"SingSong-Server/internal/pkg"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/pinecone-io/go-pinecone/pinecone"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/protobuf/types/known/structpb"
	"log"
	"net/http"
	"strconv"
)

type relatedSongResponse struct {
	SongNumber int      `json:"songNumber"`
	SongName   string   `json:"songName"`
	SingerName string   `json:"singerName"`
	Tags       []string `json:"tags"`
	IsKeep     bool     `json:"isKeep"`
	SongTempId int64    `json:"songId"`
}

var (
	defaultSize = "20"
)

// GetRelatedSong godoc
// @Summary      연관된 노래들을 조회합니다
// @Description  연관된 노래들을 조회합니다
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        songNumber path string true "노래 번호"
// @Param        size query int false "한번에 조회할 노래 개수"
// @Success      200 {object} pkg.BaseResponseStruct(data=[]relatedSongResponse) "성공"
// @Router       /songs/{songNumber}/related [get]
// @Security BearerAuth
func RelatedSong(db *sql.DB, idxConnection *pinecone.IndexConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		songNumber := c.Param("songNumber")
		if songNumber == "" {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - cannot find songNumber in path variable", nil)
			return
		}

		value, exists := c.Get("memberId")
		if !exists {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not found", nil)
			return
		}

		memberId, ok := value.(int64)
		if !ok {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - memberId not type int64", nil)
			return
		}

		sizeStr := c.DefaultQuery("size", defaultSize)
		sizeInt, err := strconv.Atoi(sizeStr)
		if err != nil || sizeInt < 0 {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid size parameter", nil)
			return
		}

		//songNumber로 벡터 디비에서 조회
		songNumberInt, err := strconv.Atoi(songNumber)
		if err != nil {
			pkg.BaseResponse(c, http.StatusBadRequest, "error - invalid songNumber", nil)
			return
		}
		filterStruct := &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"song_number": structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"$ne": structpb.NewNumberValue(float64(songNumberInt)),
					},
				}),
				"MR": structpb.NewBoolValue(false),
			},
		}

		res, err := idxConnection.QueryByVectorId(c, &pinecone.QueryByVectorIdRequest{
			VectorId:        songNumber,
			TopK:            uint32(sizeInt),
			Filter:          filterStruct,
			IncludeValues:   true,
			IncludeMetadata: true,
		})

		response := make([]relatedSongResponse, 0, sizeInt)
		for _, each := range res.Matches {
			v := each.Vector
			response = append(response, relatedSongResponse{
				SongNumber: int(v.Metadata.Fields["song_number"].GetNumberValue()),
				SongName:   v.Metadata.Fields["song_name"].GetStringValue(),
				SingerName: v.Metadata.Fields["singer_name"].GetStringValue(),
				Tags:       convertTags(v.Metadata.Fields["ssss"]),
			})
		}

		//isKeep
		all, err := mysql.KeepLists(qm.Where("memberId = ?", memberId)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		keepIds := make([]interface{}, len(all))
		for i, keep := range all {
			keepIds[i] = keep.KeepId
		}
		keepSongs, err := mysql.KeepSongs(qm.WhereIn("keepId in ?", keepIds...)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		isKeepMap := make(map[int]bool)
		for _, keepSong := range keepSongs {
			isKeepMap[keepSong.SongNumber] = true
		}

		//songTempId
		var songNumbers []interface{}
		for _, song := range response {
			songNumbers = append(songNumbers, song.SongNumber)
		}
		slice, err := mysql.SongTempInfos(qm.WhereIn("songNumber IN ?", songNumbers...)).All(c, db)
		if err != nil {
			pkg.BaseResponse(c, http.StatusInternalServerError, "error - "+err.Error(), nil)
			return
		}
		songTempIdMap := make(map[int]int64)
		for _, song := range slice {
			songTempIdMap[song.SongNumber] = song.SongTempId
		}

		// response에 isKeep과 songTempId 추가
		for i, song := range response {
			response[i].IsKeep = isKeepMap[song.SongNumber]
			response[i].SongTempId = songTempIdMap[song.SongNumber]
		}

		pkg.BaseResponse(c, http.StatusOK, "ok", response)
	}
}

func convertTags(tag *structpb.Value) []string {
	ssssField := tag.GetListValue().AsSlice()
	ssssArray := make([]string, len(ssssField))
	for i, eTag := range ssssField {
		ssssArray[i] = eTag.(string)
	}
	koreanTags, err := MapTagsEnglishToKorean(ssssArray)
	if err != nil {
		log.Printf("Failed to convert tags to korean, error: %+v", err)
		koreanTags = []string{}
	}
	return koreanTags
}
