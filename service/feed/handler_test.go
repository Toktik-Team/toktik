package main

import (
	"bou.ke/monkey"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	"toktik/repo"
	"toktik/repo/model"
	"toktik/storage"
)

const mockVideoCount = 50

var (
	testVideos = make([]*model.Video, mockVideoCount)
	respVideos = make([]*feed.Video, mockVideoCount)
)

func TestMain(m *testing.M) {
	now := int64(0) // time.Now().UnixMilli()
	for i := 0; i < mockVideoCount; i++ {
		test := &model.Video{
			Model: model.Model{
				ID:        uint32(i),
				CreatedAt: time.UnixMilli(now).Add(time.Duration(i) * time.Second),
			},
			UserId:    65535,
			Title:     "Test Video " + strconv.Itoa(i),
			FileName:  "test_video_file_" + strconv.Itoa(i) + ".mp4",
			CoverName: "test_video_cover_file_" + strconv.Itoa(i) + ".png",
		}
		resp := &feed.Video{
			Id:            uint32(i),
			Author:        &user.User{Id: 65535},
			PlayUrl:       "https://test.com/test_video_file_" + strconv.Itoa(i) + ".mp4",
			CoverUrl:      "https://test.com/test_video_cover_file_" + strconv.Itoa(i) + ".png",
			FavoriteCount: 0,     // TODO
			CommentCount:  0,     // TODO
			IsFavorite:    false, // TODO
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
	}}

	expectedNextTime := testVideos[biz.VideoCount-1].CreatedAt.Add(time.Duration(-1)).UnixMilli()
	var successResp = &feed.ListFeedResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		NextTime:   &expectedNextTime,
		Videos:     respVideos[:biz.VideoCount],
	}

	monkey.Patch(storage.GetLink, func(fileName string) (string, error) {
		return "https://test.com/" + fileName, nil
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

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	for _, v := range testVideos[:biz.VideoCount] {
		rows.AddRow(v.ID, v.CreatedAt, nil, nil, v.UserId, v.Title, v.FileName, v.CoverName)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."created_at" <= $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."created_at" DESC LIMIT ` + strconv.Itoa(biz.VideoCount))).
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
			if len(gotResp.Videos) != len(tt.wantResp.Videos) {
				t.Errorf("ListVideos() lens got %v, want %v", len(gotResp.Videos), len(tt.wantResp.Videos))
			}
			if len(gotResp.Videos) != biz.VideoCount {
				t.Errorf("ListVideos() lens got %v, want %v", len(gotResp.Videos), biz.VideoCount)
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
