package handler

//import (
//	"SingSong-Server/internal/pkg"
//	"github.com/stretchr/testify/assert"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//)
//
//func TestCreatePlaylist(t *testing.T) {
//	//t.Parallel()
//	db, mock, err := pkg.NewTestDB()
//	assert.NoError(t, err)
//	//if err != nil {
//	//	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	//}
//	////oldDB := boil.GetDB()
//	//defer func() {
//	//	db.Close()
//	//	//boil.SetDB(oldDB)
//	//}()
//	////boil.SetDB(db)
//
//	// expected arguments
//	keepName := "테스트플레이리스트"
//	memberId := int64(1)
//
//	// prepare mock expectations
//	mock.ExpectExec("INSERT INTO `keep_list` (`member_id`,`keep_name`,`deleted_at`) VALUES (?,?,?)").
//		WithArgs(memberId, keepName, nil).
//		WillReturnResult(sqlmock.NewResult(1, 1))
//
//	// execute function
//	CreatePlaylist(db, keepName, memberId)
//
//	// check if expectations were met
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
