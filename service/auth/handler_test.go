package main

import (
	"context"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
	"toktik/kitex_gen/douyin/auth"
	"toktik/service/auth/model"
)

var successArg = struct {
	ctx context.Context
	req *auth.AuthenticateRequest
}{ctx: context.Background(), req: &auth.AuthenticateRequest{Token: "authenticated-token"}}

var successResp = &auth.AuthenticateResponse{
	StatusCode: 0,
	StatusMsg:  "success",
	UserId:     114514,
}

type MockUserTokenDo struct {
	mock.Mock
}

func (m *MockUserTokenDo) First() (*model.UserToken, error) {
	//TODO implement me
	panic("implement me")
}

func TestAuthServiceImpl_Authenticate(t *testing.T) {
	testDo := new(MockUserTokenDo)
	testDo.On("First", successArg).Return(model.UserToken{
		Token:    "",
		Username: "",
		UserID:   0,
		Role:     "",
	}, nil)
	type args struct {
		ctx context.Context
		req *auth.AuthenticateRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *auth.AuthenticateResponse
		wantErr  bool
	}{
		{name: "should authenticate success", args: successArg, wantResp: successResp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthServiceImpl{}
			gotResp, err := s.Authenticate(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("Authenticate() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
