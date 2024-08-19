package handler

import (
	"SingSong-Server/internal/pkg"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSongInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 테스트 Context 생성
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 테스트 디비 생성
	db, mock, err := pkg.NewTestDB()
	assert.NoError(t, err)

	tests := []struct {
		name               string
		songInfoId         string
		memberId           interface{}
		mockQueries        func(mock sqlmock.Sqlmock)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:       "Success",
			songInfoId: "1",
			memberId:   1,
			mockQueries: func(mock sqlmock.Sqlmock) {
				// Mock the song info query
				mock.ExpectQuery("SELECT `song_info`.* FROM `song_info` WHERE (song_info_id = ?) LIMIT 1;").
					WithArgs("1").
					WillReturnRows(sqlmock.NewRows([]string{"song_info_id", "song_number", "song_name", "artist_name", "tags", "album", "octave", "is_mr"}).
						AddRow(1, 1, "Test Song", "Test Artist", "tag1,tag2", "Test Album", "Test Octave", true))

				// Mock the keep list query
				mock.ExpectQuery("SELECT `keep_list`.* FROM `keep_list` WHERE (member_id = ?);").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"keep_list_id"}).
						AddRow(1))

				// Mock the keep songs query
				mock.ExpectQuery("SELECT COUNT(*) FROM `keep_song` WHERE (`keep_list_id` IN (?)) AND (song_info_id = ?) AND (deleted_at IS NULL) LIMIT 1;").
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(1))

				// Mock the comment count query
				mock.ExpectQuery("SELECT COUNT(*) FROM `comment` WHERE (song_info_id = ? AND deleted_at is null);").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(100))

				// Mock the keep count query
				mock.ExpectQuery("SELECT COUNT(*) FROM `keep_song` WHERE (song_info_id = ? AND deleted_at is null);").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(5))
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	c.Params = gin.Params{{Key: "songId", Value: tests[0].songInfoId}}
	c.Set("memberId", tests[0].memberId)
	tests[0].mockQueries(mock)

	handler := GetSongInfo(db)
	handler(c)

	assert.Equal(t, tests[0].expectedStatusCode, w.Code)
}
