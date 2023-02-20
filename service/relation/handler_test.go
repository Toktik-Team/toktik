package main

import (
	"context"
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
		UserId: 1,
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
	defer mock.MockConn.Close()

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."user_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(1).
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
	mockUser := user.User{Id: req.UserId}
	return &user.UserResponse{StatusCode: biz.OkStatusCode, User: &mockUser}, nil
}

func TestRelationServiceImpl_GetFollowerList(t *testing.T) {
	var successArg = struct {
		ctx context.Context
		req *relation.FollowerListRequest
	}{ctx: context.Background(), req: &relation.FollowerListRequest{
		UserId: 1,
	}}

	var successResp = &relation.FollowerListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
		UserList: []*user.User{
			mockUserB,
		},
	}

	UserClient = MockUserClient{}

	relationRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationRows.AddRow(mockUserA.Id, mockUserB.Id)
	defer mock.MockConn.Close()

	mock.DBMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "relations" WHERE "relations"."target_id" = $1 AND "relations"."deleted_at" IS NULL`)).
		WithArgs(1).
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
		UserId:   mockUserA.Id,
		ToUserId: mockUserB.Id,
	}}

	var successResp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.OkStatusMsg,
	}

	var followMyselfArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		UserId:   mockUserA.Id,
		ToUserId: mockUserA.Id,
	}}

	var followMyselfResp = &relation.RelationActionResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	var followAlreadyExistsSelfArg = struct {
		ctx context.Context
		req *relation.RelationActionRequest
	}{ctx: context.Background(), req: &relation.RelationActionRequest{
		UserId:   mockUserA.Id,
		ToUserId: mockUserC.Id,
	}}

	var followAlreadyExistsSelfResp = &relation.RelationActionResponse{
		StatusCode: biz.RelationAlreadyExists,
		StatusMsg:  biz.BadRequestStatusMsg,
	}

	relationRows := sqlmock.NewRows([]string{"user_id", "target_id"})
	relationRows.AddRow(mockUserA.Id, mockUserB.Id)
	defer mock.MockConn.Close()

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
		{name: "Follow already exists", args: followAlreadyExistsSelfArg, wantResp: followAlreadyExistsSelfResp},
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
