package main

import (
	"bou.ke/monkey"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"toktik/kitex_gen/douyin/publish"
	"toktik/repo"
	"toktik/storage"
)

var testVideo []byte

func TestMain(m *testing.M) {
	var err error
	_, currentFilePath, _, _ := runtime.Caller(0)

	testVideo, err = os.ReadFile(path.Join(path.Dir(currentFilePath), "resources/bear.mp4"))
	if err != nil {
		panic("Cannot find test resources.")
	}

	code := m.Run()
	os.Exit(code)
}

func TestPublishServiceImpl_CreateVideo(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}{ctx: context.Background(), req: &publish.CreateVideoRequest{
		UserId: 1,
		Data:   testVideo,
		Title:  "Video for test",
	}}

	var successResp = &publish.CreateVideoResponse{
		Id: 1,
	}

	var invalidContentArg = struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}{ctx: context.Background(), req: &publish.CreateVideoRequest{
		UserId: 1,
		Data:   []byte{1, 2},
		Title:  "Invalid content",
	}}

	monkey.Patch(storage.Upload, func(fileName string, content io.Reader) (*s3.PutObjectOutput, error) {
		// TODO: nothing
		return nil, nil
	})

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
	defer db.Close()
	repo.SetDefault(DB)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "videos" 
    		("created_at","updated_at","deleted_at","user_id","title","file_name","cover_name") 
			VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			1,
			"Video for test",
			sqlmock.AnyArg(),
			sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	type args struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *publish.CreateVideoResponse
		wantErr  bool
	}{
		{name: "should create success", args: successArg, wantResp: successResp},
		{name: "invalid content type", args: invalidContentArg, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PublishServiceImpl{}
			gotResp, err := s.CreateVideo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("CreateVideo() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
