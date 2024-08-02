package handler

//
//import (
//	"SingSong-Server/internal/db/mysql"
//	"SingSong-Server/internal/pkg"
//	"database/sql"
//	"github.com/gin-gonic/gin"
//)
//
//// CommentOnSong godoc
//// @Summary      SongNumber에 댓글 달기
//// @Description  SongNumber에 댓글 달기
//// @Tags         Comment
//// @Accept       json
//// @Produce      json
//// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
//// @Router       /comment [post]
//// @Security BearerAuth
//func CommentOnSong(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// memberId가져오기
//		memberId, err := c.Get("memberId")
//		if !err {
//			// memberId가 없을 경우
//			pkg.BaseResponse(c, 400, "memberId not found", nil)
//			return
//		}
//		// 댓글 달기
//		//mysql.Comment{memberId: memberId, songNumber: c.PostForm("songNumber"), comment: c.PostForm("comment")}.Insert()
//
//		// 댓글 달기 성공시
//		// 댓글 정보 반환
//	}
//}
//
//// GetCommentOnSong godoc
//// @Summary      해당하는 SongNumber에 댓글을 가져옵니다
//// @Description  해당하는 SongNumber에 댓글을 가져옵니다
//// @Tags         Comment
//// @Accept       json
//// @Produce      json
//// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
//// @Router       /comment [get]
//// @Security BearerAuth
//func GetCommentOnSong(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// memberId가져오기
//
//		// 댓글 가져오기
//
//		// 댓글 가져오기 성공시
//	}
//}
//
//// ReportComment godoc
//// @Summary      해당하는 댓글ID를 통해 신고하기
//// @Description  해당하는 댓글ID를 통해 신고하기
//// @Tags         Comment
//// @Accept       json
//// @Produce      json
//// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
//// @Router       /comment/report [post]
//// @Security BearerAuth
//func ReportComment(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// memberId가져오기
//
//		// 댓글 신고하기
//	}
//}
//
//// LikeComment godoc
//// @Summary      해당하는 댓글에 좋아요 누르기
//// @Description  해당하는 댓글에 좋아요 누르기
//// @Tags         Comment
//// @Accept       json
//// @Produce      json
//// @Success      200 {object} pkg.BaseResponseStruct{data=[]PlaylistAddResponse} "성공"
//// @Router       /comment/like [post]
//// @Security BearerAuth
//func LikeComment(db *sql.DB) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// memberId가져오기
//
//		// 댓글 좋아요 누르기
//	}
//}
