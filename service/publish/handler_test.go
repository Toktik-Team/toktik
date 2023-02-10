package main

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"io"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"testing"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/publish"
	"toktik/storage"
	"toktik/test/mock"

	"bou.ke/monkey"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
		StatusCode: biz.OkStatusCode,
	}

	var invalidContentArg = struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}{ctx: context.Background(), req: &publish.CreateVideoRequest{
		UserId: 1,
		Data:   []byte{1, 2},
		Title:  "Invalid content",
	}}

	var invalidContentResp = &publish.CreateVideoResponse{
		StatusCode: biz.InvalidContentType,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	monkey.Patch(storage.Upload, func(fileName string, content io.Reader) (*s3.PutObjectOutput, error) {
		// TODO: nothing
		return nil, nil
	})
	defer mock.MockConn.Close()
	mock.DBMock.ExpectBegin()
	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "videos" 
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
	mock.DBMock.ExpectCommit()

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
		{name: "invalid content type", args: invalidContentArg, wantResp: invalidContentResp},
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
