package main

import (
	"context"
	"toktik/kitex_gen/douyin/wechat"
)

// WechatServiceImpl implements the last service interface defined in the IDL.
type WechatServiceImpl struct{}

// WechatChat implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatChat(ctx context.Context, req *wechat.MessageChatRequest) (resp *wechat.MessageChatResponse, err error) {
	// TODO: Your code here...
	return
}

// WechatAction implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatAction(ctx context.Context, req *wechat.MessageActionRequest) (resp *wechat.MessageActionResponse, err error) {
	// TODO: Your code here...
	return
}
