package main

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"toktik/constant/biz"
	relation "toktik/kitex_gen/douyin/relation"
	"toktik/kitex_gen/douyin/user"
	"toktik/test/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloudwego/kitex/client/callopt"
)

var mockUserA = &user.User{
	Id: 1,
}
var mockUserB = &user.User{
	Id: 2,
}
var mockUserC = &user.User{
	Id: 3,
}

func TestRelationServiceImpl_GetFollowList(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.FollowListRequest
	}{ctx: context.Background(), req: &relation.FollowListRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserA.Id,
	}}

	var successResp = &relation.FollowListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList: []*user.User{
			mockUserB,
		},
	}

	UserClient = MockUserClient{}

	videoRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	videoRows.AddRow(mockUserA.Id, mockUserB.Id)

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."user_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(mockUserA.Id).
		WillReturnRows(videoRows)

	type args struct {
		ctx context.Context
		req *relation.FollowListRequest
	}
	tests := []struct {
		name     string
		s        *RelationServiceImpl
		args     args
		wantResp *relation.FollowListResponse
		wantErr  bool
	}{
		{name: "GetFollowList Successful case", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationServiceImpl{}
			gotResp, err := s.GetFollowList(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationServiceImpl.GetFollowList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RelationServiceImpl.GetFollowList() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

type MockUserClient struct {
}

func (m MockUserClient) GetUser(ctx context.Context, req *user.UserRequest, callOptions ...callopt.Option) (r *user.UserResponse, err error) {
	if req.UserId == 999 {
		return &user.UserResponse{StatusCode: biz.UserNotFound, User: nil}, nil
	}
	mockUser := user.User{Id: req.UserId}
	return &user.UserResponse{StatusCode: biz.OkStatusCode, User: &mockUser}, nil
}

func TestRelationServiceImpl_GetFollowerList(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.FollowerListRequest
	}{ctx: context.Background(), req: &relation.FollowerListRequest{
		ActorId: mockUserB.Id,
		UserId:  mockUserA.Id,
	}}

	var successResp = &relation.FollowerListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList: []*user.User{
			mockUserA,
		},
	}

	UserClient = MockUserClient{}

	relationRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationRows.AddRow(mockUserA.Id, mockUserB.Id)

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."target_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(mockUserA.Id).
		WillReturnRows(relationRows)

	type args struct {
		ctx context.Context
		req *relation.FollowerListRequest
	}
	tests := []struct {
		name     string
		s        *RelationServiceImpl
		args     args
		wantResp *relation.FollowerListResponse
		wantErr  bool
	}{
		{name: "GetFollowerList Successful case", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationServiceImpl{}
			gotResp, err := s.GetFollowerList(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationServiceImpl.GetFollowerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RelationServiceImpl.GetFollowerList() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestRelationServiceImpl_Follow(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserB.Id,
	}}

	var successResp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
	}

	var followMyselfArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserA.Id,
	}}

	var followMyselfResp = &relation.RelationActionResponse{
		StatusCode: biz.InvalidToUserId,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	var followAlreadyExistsArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserC.Id,
	}}

	var followAlreadyExistsResp = &relation.RelationActionResponse{
		StatusCode: biz.RelationAlreadyExists,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	var followUserNotFoundArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  999,
	}}

	var followUserNotFoundResp = &relation.RelationActionResponse{
		StatusCode: biz.UserNotFound,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	UserClient = MockUserClient{}

	mock.DBMock.ExpectBegin()

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "relations" ("created_at","updated_at","deleted_at","user_id","target_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			mockUserA.Id,
			mockUserB.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.DBMock.ExpectCommit()
	mock.DBMock.ExpectBegin()

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "relations" ("created_at","updated_at","deleted_at","user_id","target_id") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
		WithArgs(sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			mockUserA.Id,
			mockUserC.Id).
		WillReturnError(fmt.Errorf("Duplicated row"))

	mock.DBMock.ExpectRollback()

	type args struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}
	tests := []struct {
		name     string
		s        *RelationServiceImpl
		args     args
		wantResp *relation.RelationActionResponse
		wantErr  bool
	}{
		{name: "Follow Successful case", args: successArg, wantResp: successResp},
		{name: "Follow myself", args: followMyselfArg, wantResp: followMyselfResp},
		{name: "Follow already exists", args: followAlreadyExistsArg, wantResp: followAlreadyExistsResp, wantErr: true},
		{name: "Follow user not found", args: followUserNotFoundArg, wantResp: followUserNotFoundResp, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationServiceImpl{}
			gotResp, err := s.Follow(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationServiceImpl.Follow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RelationServiceImpl.Follow() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestRelationServiceImpl_Unfollow(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserB.Id,
	}}

	var successResp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
	}

	var unfollowNotFoundArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserC.Id,
	}}

	var unfollowNotFoundResp = &relation.RelationActionResponse{
		StatusCode: biz.RelationNotFound,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	UserClient = MockUserClient{}

	relationRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationRows.AddRow(mockUserA.Id, mockUserB.Id)

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."user_id" = $1 AND "relations"."target_id" = $2 AND "relations"."deleted_at" IS NULL ORDER BY "relations"."id" LIMIT 1`)).
		WithArgs(mockUserA.Id, mockUserB.Id).
		WillReturnRows(relationRows)

	mock.DBMock.ExpectBegin()

	mock.DBMock.ExpectExec(regexp.QuoteMeta(`UPDATE "relations" SET "deleted_at"=$1 WHERE "relations"."user_id" = $2 AND "relations"."target_id" = $3 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.DBMock.ExpectCommit()

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."user_id" = $1 AND "relations"."target_id" = $2 AND "relations"."deleted_at" IS NULL ORDER BY "relations"."id" LIMIT 1`)).
		WithArgs(mockUserA.Id, mockUserC.Id).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "target_id"}))

	type args struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}
	tests := []struct {
		name     string
		s        *RelationServiceImpl
		args     args
		wantResp *relation.RelationActionResponse
		wantErr  bool
	}{
		{name: "Unfollow Successful case", args: successArg, wantResp: successResp},
		{name: "Unfollow not found", args: unfollowNotFoundArg, wantResp: unfollowNotFoundResp, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationServiceImpl{}
			gotResp, err := s.Unfollow(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationServiceImpl.Unfollow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RelationServiceImpl.Follow() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestRelationServiceImpl_GetFriendList(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.FriendListRequest
	}{ctx: context.Background(), req: &relation.FriendListRequest{
		ActorId: mockUserA.Id,
		UserId:  mockUserA.Id,
	}}

	var successResp = &relation.FriendListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList: []*user.User{
			biz.ChatGPTUser,
			mockUserB,
		},
	}

	UserClient = MockUserClient{}

	relationFollowerRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationFollowerRows.AddRow(mockUserB.Id, mockUserA.Id)

	relationFollowRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationFollowRows.AddRow(mockUserA.Id, mockUserB.Id)

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."target_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(mockUserA.Id).
		WillReturnRows(relationFollowerRows)

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."user_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(mockUserA.Id).
		WillReturnRows(relationFollowRows)

	type args struct {
		ctx context.Context
		req *relation.FriendListRequest
	}
	tests := []struct {
		name     string
		s        *RelationServiceImpl
		args     args
		wantResp *relation.FriendListResponse
		wantErr  bool
	}{
		{name: "GetFriends Successful case", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RelationServiceImpl{}
			gotResp, err := s.GetFriendList(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RelationServiceImpl.GetFriendList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("RelationServiceImpl.GetFriendList() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
