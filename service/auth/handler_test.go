package main

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"testing"
	"toktik/kitex_gen/douyin/auth"
	"toktik/repo"
)

var successArg = struct {
	ctx context.Context
	req *auth.AuthenticateRequest
}{ctx: context.Background(), req: &auth.AuthenticateRequest{Token: "authenticated-token"}}

var successResp = &auth.AuthenticateResponse{
	StatusCode: 0,
	StatusMsg:  "success",
	UserId:     114514,
}

func TestAuthServiceImpl_Authenticate(t *testing.T) {
	mock, db := NewDBMock(t)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "user_tokens" WHERE "user_tokens"."token" = $1 ORDER BY "user_tokens"."token" LIMIT 1`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "token", "created_at", "updated_at"}).
			AddRow(1, 114514, "authenticated-token", "2021-01-01 00:00:00", "2021-01-01 00:00:00"))

	type args struct {
		ctx context.Context
		req *auth.AuthenticateRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.AuthenticateResponse
		wantErr  bool
	}{
		{name: "should authenticate success", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServiceImpl{}
			gotResp, err := s.Authenticate(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Authenticate() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func NewDBMock(t *testing.T) (sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	DB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo.SetDefault(DB)
	return mock, db
}
