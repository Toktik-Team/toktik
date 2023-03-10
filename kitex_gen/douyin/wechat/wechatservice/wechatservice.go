// Code generated by Kitex v0.4.4. DO NOT EDIT.

package wechatservice

import (
	"context"
	"fmt"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	wechat "toktik/kitex_gen/douyin/wechat"
)

func serviceInfo() *kitex.ServiceInfo {
	return wechatServiceServiceInfo
}

var wechatServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "WechatService"
	handlerType := (*wechat.WechatService)(nil)
	methods := map[string]kitex.MethodInfo{
		"WechatChat":   kitex.NewMethodInfo(wechatChatHandler, newWechatChatArgs, newWechatChatResult, false),
		"WechatAction": kitex.NewMethodInfo(wechatActionHandler, newWechatActionArgs, newWechatActionResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "douyin.wechat",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func wechatChatHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(wechat.MessageChatRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(wechat.WechatService).WechatChat(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *WechatChatArgs:
		success, err := handler.(wechat.WechatService).WechatChat(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*WechatChatResult)
		realResult.Success = success
	}
	return nil
}
func newWechatChatArgs() interface{} {
	return &WechatChatArgs{}
}

func newWechatChatResult() interface{} {
	return &WechatChatResult{}
}

type WechatChatArgs struct {
	Req *wechat.MessageChatRequest
}

func (p *WechatChatArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(wechat.MessageChatRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *WechatChatArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *WechatChatArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *WechatChatArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in WechatChatArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *WechatChatArgs) Unmarshal(in []byte) error {
	msg := new(wechat.MessageChatRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var WechatChatArgs_Req_DEFAULT *wechat.MessageChatRequest

func (p *WechatChatArgs) GetReq() *wechat.MessageChatRequest {
	if !p.IsSetReq() {
		return WechatChatArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *WechatChatArgs) IsSetReq() bool {
	return p.Req != nil
}

type WechatChatResult struct {
	Success *wechat.MessageChatResponse
}

var WechatChatResult_Success_DEFAULT *wechat.MessageChatResponse

func (p *WechatChatResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(wechat.MessageChatResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *WechatChatResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *WechatChatResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *WechatChatResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in WechatChatResult")
	}
	return proto.Marshal(p.Success)
}

func (p *WechatChatResult) Unmarshal(in []byte) error {
	msg := new(wechat.MessageChatResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *WechatChatResult) GetSuccess() *wechat.MessageChatResponse {
	if !p.IsSetSuccess() {
		return WechatChatResult_Success_DEFAULT
	}
	return p.Success
}

func (p *WechatChatResult) SetSuccess(x interface{}) {
	p.Success = x.(*wechat.MessageChatResponse)
}

func (p *WechatChatResult) IsSetSuccess() bool {
	return p.Success != nil
}

func wechatActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(wechat.MessageActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(wechat.WechatService).WechatAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *WechatActionArgs:
		success, err := handler.(wechat.WechatService).WechatAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*WechatActionResult)
		realResult.Success = success
	}
	return nil
}
func newWechatActionArgs() interface{} {
	return &WechatActionArgs{}
}

func newWechatActionResult() interface{} {
	return &WechatActionResult{}
}

type WechatActionArgs struct {
	Req *wechat.MessageActionRequest
}

func (p *WechatActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(wechat.MessageActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *WechatActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *WechatActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *WechatActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in WechatActionArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *WechatActionArgs) Unmarshal(in []byte) error {
	msg := new(wechat.MessageActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var WechatActionArgs_Req_DEFAULT *wechat.MessageActionRequest

func (p *WechatActionArgs) GetReq() *wechat.MessageActionRequest {
	if !p.IsSetReq() {
		return WechatActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *WechatActionArgs) IsSetReq() bool {
	return p.Req != nil
}

type WechatActionResult struct {
	Success *wechat.MessageActionResponse
}

var WechatActionResult_Success_DEFAULT *wechat.MessageActionResponse

func (p *WechatActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(wechat.MessageActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *WechatActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *WechatActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *WechatActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in WechatActionResult")
	}
	return proto.Marshal(p.Success)
}

func (p *WechatActionResult) Unmarshal(in []byte) error {
	msg := new(wechat.MessageActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *WechatActionResult) GetSuccess() *wechat.MessageActionResponse {
	if !p.IsSetSuccess() {
		return WechatActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *WechatActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*wechat.MessageActionResponse)
}

func (p *WechatActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) WechatChat(ctx context.Context, Req *wechat.MessageChatRequest) (r *wechat.MessageChatResponse, err error) {
	var _args WechatChatArgs
	_args.Req = Req
	var _result WechatChatResult
	if err = p.c.Call(ctx, "WechatChat", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) WechatAction(ctx context.Context, Req *wechat.MessageActionRequest) (r *wechat.MessageActionResponse, err error) {
	var _args WechatActionArgs
	_args.Req = Req
	var _result WechatActionResult
	if err = p.c.Call(ctx, "WechatAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
