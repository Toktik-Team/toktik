package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"toktik/kitex_gen/douyin/wechat"
)

func TestWechatServiceImpl_WechatAction(t *testing.T) {
	type args struct {
		ctx context.Context
		req *wechat.MessageActionRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *wechat.MessageActionResponse
		wantErr  bool
	}{
		// TODO: Use testcontainers to test redis
		//{"should pass", args{context.Background(), &wechat.MessageActionRequest{SenderId: 114, ReceiverId: 514, ActionType: 1, Content: `Hi, nice 2 meet u.`}}, &wechat.MessageActionResponse{StatusCode: 0, StatusMsg: "success"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WechatServiceImpl{}
			gotResp, err := s.WechatAction(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("WechatAction() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestWechatServiceImpl_WechatChat(t *testing.T) {
	type args struct {
		ctx context.Context
		req *wechat.MessageChatRequest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *wechat.MessageChatResponse
		wantErr  bool
	}{
		// TODO: Use testcontainers to test redis
		//{"should pass", args{context.Background(), &wechat.MessageChatRequest{SenderId: 114, ReceiverId: 514}}, &wechat.MessageChatResponse{StatusCode: 0, StatusMsg: "success"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WechatServiceImpl{}
			gotResp, err := s.WechatChat(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("WechatChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("WechatChat() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func TestWechatServiceImpl_generateKey(t *testing.T) {
	test1Result := "chat:0:0"
	test2Result := "chat:1:2"
	test3Result := "chat:1:2"
	test4Result := fmt.Sprintf("chat:%d:%d", 1<<32-1, 1<<32-1)
	test5Result := fmt.Sprintf("chat:%d:%d", 1<<32-2, 1<<32-1)

	getUint32Ptr := func(i uint32) *uint32 {
		return &i
	}

	type args struct {
		uid1 *uint32
		uid2 *uint32
	}
	tests := []struct {
		name string
		args args
		want func(s *string) bool
	}{
		{"test1", args{uid1: new(uint32), uid2: new(uint32)}, func(s *string) bool {
			return *s == test1Result
		}},
		{"test2", args{uid1: getUint32Ptr(1), uid2: getUint32Ptr(2)}, func(s *string) bool {
			return *s == test2Result
		}},
		{"test3", args{uid1: getUint32Ptr(2), uid2: getUint32Ptr(1)}, func(s *string) bool {
			return *s == test3Result
		}},
		{"test4", args{uid1: getUint32Ptr(1<<32 - 1), uid2: getUint32Ptr(1<<32 - 1)}, func(s *string) bool {
			return *s == test4Result
		}},
		{"test5", args{uid1: getUint32Ptr(1<<32 - 2), uid2: getUint32Ptr(1<<32 - 1)}, func(s *string) bool {
			return *s == test5Result
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &WechatServiceImpl{}
			if got := s.generateKey(tt.args.uid1, tt.args.uid2); !tt.want(got) {
				t.Errorf("generateKey() = %v", got)
			}
		})
	}
}
