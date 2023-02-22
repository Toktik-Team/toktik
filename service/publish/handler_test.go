package main

import (
	"context"
	"database/sql"
	"io"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"testing"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/publish"
	"toktik/kitex_gen/douyin/user"
	"toktik/repo/model"
	"toktik/storage"
	"toktik/test/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudwego/kitex/client/callopt"
)

var testVideo []byte

var mockUser = user.User{Id: 65535}
var (
	mockVideoReq = model.Video{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.UnixMilli(0),
		},
		UserId:    mockUser.Id,
		Title:     "Test Video " + strconv.Itoa(1),
		FileName:  "test_video_file_" + strconv.Itoa(1) + ".mp4",
		CoverName: "test_video_cover_file_" + strconv.Itoa(1) + ".png",
	}
	mockVideoResp = feed.Video{
		Id:            1,
		Author:        &mockUser,
		PlayUrl:       "https://test.com/test_video_file_" + strconv.Itoa(1) + ".mp4",
		CoverUrl:      "https://test.com/test_video_cover_file_" + strconv.Itoa(1) + ".png",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         "Test Video " + strconv.Itoa(1),
	}
)

func TestMain(m *testing.M) {
	var err error
	_, currentFilePath, _, _ := runtime.Caller(0)

	testVideo, err = os.ReadFile(path.Join(path.Dir(currentFilePath), "resources/bear.mp4"))
	if err != nil {
		panic("Cannot find test resources.")
	}

	storage.Instance = MockStorageProvider{}

	code := m.Run()
	os.Exit(code)
}

func TestPublishServiceImpl_CreateVideo(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}{ctx: context.Background(), req: &publish.CreateVideoRequest{
		ActorId: 1,
		Data:    testVideo,
		Title:   "Video for test",
	}}

	var successResp = &publish.CreateVideoResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.PublishActionSuccess,
	}

	var invalidContentArg = struct {
		ctx context.Context
		req *publish.CreateVideoRequest
	}{ctx: context.Background(), req: &publish.CreateVideoRequest{
		ActorId: 1,
		Data:    []byte{1, 2},
		Title:   "Invalid content",
	}}

	var invalidContentResp = &publish.CreateVideoResponse{
		StatusCode: biz.InvalidContentType,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	defer func(MockConn *sql.DB) {
		err := MockConn.Close()
		if err != nil {
			panic(err)
		}
	}(mock.Conn)
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

func TestPublishServiceImpl_ListVideo(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *publish.ListVideoRequest
	}{ctx: context.Background(), req: &publish.ListVideoRequest{
		ActorId: 1,
		UserId:  65535,
	}}

	var successResp = &publish.ListVideoResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  []*feed.Video{&mockVideoResp},
	}

	FeedClient = MockFeedClient{}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	rows.AddRow(mockVideoReq.ID, mockVideoReq.CreatedAt, nil, nil, mockVideoReq.UserId, mockVideoReq.Title, mockVideoReq.FileName, mockVideoReq.CoverName)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."user_id" = $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."created_at" DESC`)).
		WillReturnRows(rows)

	type args struct {
		ctx context.Context
		req *publish.ListVideoRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *publish.ListVideoResponse
		wantErr  bool
	}{
		{name: "List Video", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PublishServiceImpl{}
			gotResp, err := s.ListVideo(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ListVideo() gotResp %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

type MockStorageProvider struct {
}

func (m MockStorageProvider) Upload(string, io.Reader) (*storage.PutObjectOutput, error) {
	// Nothing to do
	return &storage.PutObjectOutput{}, nil
}

func (m MockStorageProvider) GetLink(fileName string) (string, error) {
	return "https://test.com/" + fileName, nil
}

type MockFeedClient struct {
}

func (m MockFeedClient) ListVideos(context.Context, *feed.ListFeedRequest, ...callopt.Option) (r *feed.ListFeedResponse, err error) {
	return &feed.ListFeedResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  []*feed.Video{&mockVideoResp},
	}, nil
}

func (m MockFeedClient) QueryVideos(context.Context, *feed.QueryVideosRequest, ...callopt.Option) (r *feed.QueryVideosResponse, err error) {
	return &feed.QueryVideosResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  []*feed.Video{&mockVideoResp},
	}, nil
}
