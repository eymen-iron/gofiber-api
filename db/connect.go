package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func OpenDBConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:pass@tcp(localhost:3306)/db-name")
	if err != nil {
		return nil, fmt.Errorf("Veritabanına bağlanılamadı: %v", err)
	}
	return db, nil
}
