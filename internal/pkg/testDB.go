package pkg

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
)

func NewTestDB() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error {
		if actualSQL != expectedSQL {
			fmt.Printf("Mismatch in SQL query detected\nActual query: %s\nExpected query: %s\n", actualSQL, expectedSQL)
			return fmt.Errorf("Mismatch in SQL query detected")
		}
		return nil
	})))
}
