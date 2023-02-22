package main

import (
	"context"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/favorite"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	"toktik/repo/model"
	"toktik/storage"
	"toktik/test/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudwego/kitex/client/callopt"
)

const mockVideoCount = 50

var (
	testVideos = make([]*model.Video, mockVideoCount)
	respVideos = make([]*feed.Video, mockVideoCount)
)

var mockUser = user.User{Id: 65535}

func TestMain(m *testing.M) {
	storage.Instance = MockStorageProvider{}

	now := time.Now().UnixMilli()
	for i := 0; i < mockVideoCount; i++ {
		test := &model.Video{
			Model: model.Model{
				ID:        uint32(i),
				CreatedAt: time.UnixMilli(now).Add(time.Duration(i) * time.Second),
			},
			UserId:    mockUser.Id,
			Title:     "Test Video " + strconv.Itoa(i),
			FileName:  "test_video_file_" + strconv.Itoa(i) + ".mp4",
			CoverName: "test_video_cover_file_" + strconv.Itoa(i) + ".png",
		}
		resp := &feed.Video{
			Id:            uint32(i),
			Author:        &mockUser,
			PlayUrl:       "https://test.com/test_video_file_" + strconv.Itoa(i) + ".mp4",
			CoverUrl:      "https://test.com/test_video_cover_file_" + strconv.Itoa(i) + ".png",
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         "Test Video " + strconv.Itoa(i),
		}
		testVideos[i] = test
		respVideos[i] = resp
	}
	testVideos = reverseModelVideo(testVideos)
	respVideos = reverseFeedVideo(respVideos)

	code := m.Run()
	os.Exit(code)
}

func TestFeedServiceImpl_ListVideos(t *testing.T) {
	pTime := strconv.FormatInt(time.Now().Add(time.Duration(1)*time.Hour).UnixMilli(), 10)
	var successArg = struct {
		ctx context.Context
		req *feed.ListFeedRequest
	}{ctx: context.Background(), req: &feed.ListFeedRequest{
		LatestTime: &pTime,
		ActorId:    nil,
	}}

	expectedNextTime := testVideos[biz.VideoCount-1].CreatedAt.Add(time.Duration(-1)).UnixMilli()
	var successResp = &feed.ListFeedResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		NextTime:   &expectedNextTime,
		VideoList:  respVideos[:biz.VideoCount],
	}

	UserClient = MockUserClient{}
	CommentClient = MockCommentClient{}
	FavoriteClient = MockFavoriteClient{}

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	for _, v := range testVideos[:biz.VideoCount] {
		rows.AddRow(v.ID, v.CreatedAt, nil, nil, v.UserId, v.Title, v.FileName, v.CoverName)
	}

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."created_at" <= $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."created_at" DESC LIMIT ` + strconv.Itoa(biz.VideoCount))).
		WillReturnRows(rows)

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
			if len(gotResp.VideoList) != len(tt.wantResp.VideoList) {
				t.Errorf("ListVideos() lens got %v, want %v", len(gotResp.VideoList), len(tt.wantResp.VideoList))
			}
			if len(gotResp.VideoList) != biz.VideoCount {
				t.Errorf("ListVideos() lens got %v, want %v", len(gotResp.VideoList), biz.VideoCount)
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ListVideos() gotResp %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func reverseModelVideo(s []*model.Video) []*model.Video {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func reverseFeedVideo(s []*feed.Video) []*feed.Video {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

type MockUserClient struct {
}

func (m MockUserClient) GetUser(context.Context, *user.UserRequest, ...callopt.Option) (r *user.UserResponse, err error) {
	return &user.UserResponse{StatusCode: biz.OkStatusCode, User: &mockUser}, nil
}

type MockCommentClient struct {
}

func (m MockCommentClient) CountComment(context.Context, *comment.CountCommentRequest, ...callopt.Option) (r *comment.CountCommentResponse, err error) {
	return &comment.CountCommentResponse{
		StatusCode:   biz.OkStatusCode,
		StatusMsg:    &biz.OkStatusMsg,
		CommentCount: 0,
	}, nil
}

func (m MockCommentClient) ActionComment(context.Context, *comment.ActionCommentRequest, ...callopt.Option) (r *comment.ActionCommentResponse, err error) {
	panic("unimplemented")
}

func (m MockCommentClient) ListComment(context.Context, *comment.ListCommentRequest, ...callopt.Option) (r *comment.ListCommentResponse, err error) {
	panic("unimplemented")
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

type MockFavoriteClient struct {
}

func (m MockFavoriteClient) CountUserFavorite(context.Context, *favorite.CountUserFavoriteRequest, ...callopt.Option) (r *favorite.CountUserFavoriteResponse, err error) {
	return &favorite.CountUserFavoriteResponse{
		Count: 0,
	}, nil
}
func (m MockFavoriteClient) CountUserTotalFavorited(context.Context, *favorite.CountUserTotalFavoritedRequest, ...callopt.Option) (r *favorite.CountUserTotalFavoritedResponse, err error) {
	return &favorite.CountUserTotalFavoritedResponse{
		Count: 0,
	}, nil
}

func (m MockFavoriteClient) FavoriteAction(context.Context, *favorite.FavoriteRequest, ...callopt.Option) (r *favorite.FavoriteResponse, err error) {
	panic("unimplemented")
}

func (m MockFavoriteClient) FavoriteList(context.Context, *favorite.FavoriteListRequest, ...callopt.Option) (r *favorite.FavoriteListResponse, err error) {
	panic("unimplemented")
}

func (m MockFavoriteClient) IsFavorite(context.Context, *favorite.IsFavoriteRequest, ...callopt.Option) (r *favorite.IsFavoriteResponse, err error) {
	return &favorite.IsFavoriteResponse{
		Result: false,
	}, nil
}

func (m MockFavoriteClient) CountFavorite(context.Context, *favorite.CountFavoriteRequest, ...callopt.Option) (r *favorite.CountFavoriteResponse, err error) {
	return &favorite.CountFavoriteResponse{
		Count: 0,
	}, nil
}
