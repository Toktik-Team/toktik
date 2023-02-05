package main

import (
	"context"
	"os"
	"path"
	"reflect"
	"runtime"
	"testing"
	"toktik/kitex_gen/douyin/publish"
	"toktik/service/publish/storage"
)

func TestPublishServiceImpl_CreateVideo(t *testing.T) {
	_, currentFilePath, _, _ := runtime.Caller(0)

	testVideo, err := os.ReadFile(path.Join(path.Dir(currentFilePath), "resources/bear.mp4"))
	if err != nil {
		panic("Cannot find test resources.")
	}
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

	storage.Init()

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
		// TODO: Add test cases.
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
