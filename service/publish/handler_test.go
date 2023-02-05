package main

import (
	"bou.ke/monkey"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
	"toktik/kitex_gen/douyin/publish"
	gen "toktik/repo"
	"toktik/repo/model"
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

	var v = gen.Q.Video
	monkey.PatchInstanceMethod(reflect.TypeOf(v), "Create", func(do gen.IVideoDo, values ...*model.Video) error {
		values[0].ID = 1
		return nil
	})
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
