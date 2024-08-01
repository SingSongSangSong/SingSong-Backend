package handler

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func TestCreatePlaylist(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	oldDB := boil.GetDB()
	defer func() {
		db.Close()
		boil.SetDB(oldDB)
	}()
	boil.SetDB(db)

	// expected arguments
	keepName := "테스트플레이리스트"
	memberId := int64(1)

	// prepare mock expectations
	query := regexp.QuoteMeta("INSERT INTO `keepList` (`memberId`,`keepName`) VALUES (?,?)")
	mock.ExpectExec(query).
		WithArgs(memberId, keepName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// execute function
	CreatePlaylist(db, keepName, memberId)

	// check if expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
