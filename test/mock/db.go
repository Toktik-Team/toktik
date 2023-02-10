package mock

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"toktik/repo"
)

var DBMock sqlmock.Sqlmock
var MockConn *sql.DB

func init() {
	var err error
	MockConn, DBMock, err = sqlmock.New()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 MockConn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo.SetDefault(db)
}
