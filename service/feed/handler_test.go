package main

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
	"toktik/kitex_gen/douyin/feed"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestFeedServiceImpl_ListVideos(t *testing.T) {
	pTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	var successArg = struct {
		ctx context.Context
		req *feed.ListFeedRequest
	}{ctx: context.Background(), req: &feed.ListFeedRequest{
		LatestTime: &pTime,
	}}

	var successResp = &feed.ListFeedResponse{
		StatusCode: 0,
		//TODO
		Videos: []*feed.Video{},
	}

	type args struct {
		ctx context.Context
		req *feed.ListFeedRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *feed.ListFeedResponse
		wantErr  bool
	}{
		{name: "Feed Video", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FeedServiceImpl{}
			gotResp, err := s.ListVideos(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ListVideos() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
