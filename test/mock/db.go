package mock

import (
	"database/sql"
	"log"
	"toktik/repo"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBMock sqlmock.Sqlmock
var Conn *sql.DB

func init() {
	var err error
	Conn, DBMock, err = sqlmock.New()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 Conn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo.SetDefault(db)
}
