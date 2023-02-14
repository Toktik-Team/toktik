package main

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudwego/kitex/client/callopt"
	"reflect"
	"regexp"
	"testing"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/user"
	"toktik/repo/model"
	"toktik/test/mock"
)

var (
	testComment = &model.Comment{
		Model: model.Model{
			CreatedAt: time.Now(),
		},
		CommentId: 1,
		VideoId:   1,
		UserId:    mockUser.Id,
		Content:   "test comment",
	}
)

var (
	mockUser  = user.User{Id: 65535}
	mockVideo = model.Video{
		Model: model.Model{
			ID:        1,
			CreatedAt: time.Now(),
		},
		UserId:    mockUser.Id,
		Title:     "test video",
		FileName:  "test.mp4",
		CoverName: "test.jpg",
	}
)

func TestCommentServiceImpl_ActionComment_Add(t *testing.T) {
	var addSuccessArg = struct {
		ctx context.Context
		req *comment.ActionCommentRequest
	}{ctx: context.Background(), req: &comment.ActionCommentRequest{
		ActorId:    mockUser.Id,
		VideoId:    mockVideo.ID,
		ActionType: comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD,
		Action:     &comment.ActionCommentRequest_CommentText{CommentText: testComment.Content},
	}}
	var addSuccessResp = &comment.ActionCommentResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Comment: &comment.Comment{
			Id:         testComment.CommentId,
			User:       &mockUser,
			Content:    testComment.Content,
			CreateDate: testComment.CreatedAt.Format("01-02"),
		},
	}

	UserClient = MockUserClient{}

	videoRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	videoRows.AddRow(mockVideo.ID, mockVideo.CreatedAt, mockVideo.UpdatedAt, mockVideo.DeletedAt, mockVideo.UserId, mockVideo.Title, mockVideo.FileName, mockVideo.CoverName)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."id" = $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."id" LIMIT 1`)).
		WillReturnRows(videoRows)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "comments" WHERE "comments"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.DBMock.ExpectBegin()
	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments" ("created_at","updated_at","deleted_at","comment_id","video_id","user_id","content") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(testComment.CommentId))
	mock.DBMock.ExpectCommit()

	type args struct {
		ctx context.Context
		req *comment.ActionCommentRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *comment.ActionCommentResponse
		wantErr  bool
	}{
		{name: "Add Comment", args: addSuccessArg, wantResp: addSuccessResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentServiceImpl{}
			gotResp, err := s.ActionComment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ActionComment(Add) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ActionComment(Add) gotResp %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestCommentServiceImpl_ActionComment_Delete(t *testing.T) {
	var deleteSuccessArg = struct {
		ctx context.Context
		req *comment.ActionCommentRequest
	}{ctx: context.Background(), req: &comment.ActionCommentRequest{
		ActorId:    mockUser.Id,
		VideoId:    mockVideo.ID,
		ActionType: comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE,
		Action:     &comment.ActionCommentRequest_CommentId{CommentId: testComment.CommentId},
	}}
	var deletedSuccessResp = &comment.ActionCommentResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Comment:    nil,
	}

	UserClient = MockUserClient{}

	videoRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	videoRows.AddRow(mockVideo.ID, mockVideo.CreatedAt, mockVideo.UpdatedAt, mockVideo.DeletedAt, mockVideo.UserId, mockVideo.Title, mockVideo.FileName, mockVideo.CoverName)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."id" = $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."id" LIMIT 1`)).
		WillReturnRows(videoRows)

	commentRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "comment_id", "video_id", "user_id", "content"})
	commentRows.AddRow(testComment.ID, testComment.CreatedAt, testComment.UpdatedAt, testComment.DeletedAt, testComment.CommentId, testComment.VideoId, testComment.UserId, testComment.Content)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."video_id" = $1 AND "comments"."comment_id" = $2 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."id" LIMIT 1`)).
		WillReturnRows(commentRows)

	mock.DBMock.ExpectBegin()
	mock.DBMock.
		ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET "deleted_at"=$1 WHERE "comments"."id" = $2 AND "comments"."deleted_at" IS NULL`)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.DBMock.ExpectCommit()

	type args struct {
		ctx context.Context
		req *comment.ActionCommentRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *comment.ActionCommentResponse
		wantErr  bool
	}{
		{name: "Delete Comment", args: deleteSuccessArg, wantResp: deletedSuccessResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentServiceImpl{}
			gotResp, err := s.ActionComment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ActionComment(Delete) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ActionComment(Delete) gotResp %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestCommentServiceImpl_ListComment(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *comment.ListCommentRequest
	}{ctx: context.Background(), req: &comment.ListCommentRequest{
		ActorId: mockUser.Id,
		VideoId: mockVideo.ID,
	}}

	comments := make([]*comment.Comment, 0)
	comments = append(comments, &comment.Comment{
		Id: testComment.CommentId,
		User: &user.User{
			Id:            mockUser.Id,
			Name:          mockUser.Name,
			FollowCount:   mockUser.FollowCount,
			FollowerCount: mockUser.FollowerCount,
			IsFollow:      mockUser.IsFollow,
		},
		Content:    testComment.Content,
		CreateDate: testComment.CreatedAt.Format("01-02"),
	})
	var successResp = &comment.ListCommentResponse{
		StatusCode:  biz.OkStatusCode,
		StatusMsg:   &biz.OkStatusMsg,
		CommentList: comments,
	}

	UserClient = MockUserClient{}

	videoRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "title", "file_name", "cover_name"})
	videoRows.AddRow(mockVideo.ID, mockVideo.CreatedAt, mockVideo.UpdatedAt, mockVideo.DeletedAt, mockVideo.UserId, mockVideo.Title, mockVideo.FileName, mockVideo.CoverName)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "videos" WHERE "videos"."id" = $1 AND "videos"."deleted_at" IS NULL ORDER BY "videos"."id" LIMIT 1`)).
		WillReturnRows(videoRows)

	commentRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "comment_id", "video_id", "user_id", "content"})
	commentRows.AddRow(testComment.ID, testComment.CreatedAt, testComment.UpdatedAt, testComment.DeletedAt, testComment.CommentId, testComment.VideoId, testComment.UserId, testComment.Content)

	mock.DBMock.
		ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE "comments"."video_id" = $1 AND "comments"."deleted_at" IS NULL ORDER BY "comments"."created_at" DESC`)).
		WillReturnRows(commentRows)

	type args struct {
		ctx context.Context
		req *comment.ListCommentRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *comment.ListCommentResponse
		wantErr  bool
	}{
		{name: "List Comment", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CommentServiceImpl{}
			gotResp, err := s.ListComment(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("ListComment() gotResp %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

type MockUserClient struct {
}

func (m MockUserClient) GetUser(ctx context.Context, Req *user.UserRequest, callOptions ...callopt.Option) (r *user.UserResponse, err error) {
	return &user.UserResponse{StatusCode: biz.OkStatusCode, User: &mockUser}, nil
}
