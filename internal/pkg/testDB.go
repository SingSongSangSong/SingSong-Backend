package pkg

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

// TestDBMock 생성
func NewTestDBMock() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New()
}

// gin Context + recorder 생성
func NewTestGinContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	return c, rec
}

// SELECT COUNT(*) 쿼리 mock
func ExpectCountQuery(mock sqlmock.Sqlmock, table, where string, args []any, count int) {
	query := regexp.QuoteMeta("SELECT COUNT(*) FROM `" + table + "` WHERE " + where + " LIMIT 1;")
	mock.ExpectQuery(query).
		WithArgs(toDriverArgs(args)...).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count)).
		RowsWillBeClosed()
}

// INSERT INTO table
func ExpectInsert(mock sqlmock.Sqlmock, table string, args []any) {
	mock.ExpectExec("INSERT INTO `" + table + "`.*").
		WithArgs(toDriverArgs(args)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
}

// SELECT * FROM table WHERE id = ?
func ExpectSelectByID(mock sqlmock.Sqlmock, table, idColumn string, id any, columns []string) {
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT `" + columns[0] + "`,`" + columns[1] + "`,`" + columns[2] + "` FROM `" + table + "` WHERE `" + idColumn + "`=?",
	)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows(columns).AddRow(id, now, now)).
		RowsWillBeClosed()
}

// args → []driver.Value 변환 유틸
func toDriverArgs(args []any) []driver.Value {
	var out []driver.Value
	for _, a := range args {
		out = append(out, a)
	}
	return out
}
