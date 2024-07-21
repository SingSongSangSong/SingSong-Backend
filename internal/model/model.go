package model

import (
	"SingSong-Backend/config"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Model struct {
	db *sql.DB
}

func NewModel(config *config.MysqlConfig) (*Model, error) {
	r := &Model{}
	var err error

	// MySQL 연결 문자열 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username, config.Password, config.Host, config.Port, config.Schema)

	// MySQL 연결
	if r.db, err = sql.Open("mysql", dsn); err != nil {
		return nil, err
	}

	// 연결 확인
	if err := r.db.Ping(); err != nil {
		return nil, err
	}

	// 테이블 생성
	if err := r.createTables(); err != nil {
		return nil, err
	}

	return r, nil
}

// 테이블 생성 함수
func (model *Model) createTables() error {
	tableCreationQueries := []string{
		"CREATE TABLE IF NOT EXISTS `user` (id BIGINT AUTO_INCREMENT PRIMARY KEY, username VARCHAR(255) NOT NULL);",
	}

	for _, query := range tableCreationQueries {
		if _, err := model.db.Exec(query); err != nil {
			return fmt.Errorf("failed to create tables: %v", err)
		}
	}

	return nil
}
