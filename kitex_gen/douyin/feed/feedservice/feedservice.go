// Code generated by Kitex v0.4.4. DO NOT EDIT.

package feedservice

import (
	"context"
	"fmt"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
	feed "toktik/kitex_gen/douyin/feed"
)

func serviceInfo() *kitex.ServiceInfo {
	return feedServiceServiceInfo
}

var feedServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "FeedService"
	handlerType := (*feed.FeedService)(nil)
	methods := map[string]kitex.MethodInfo{
		"ListVideos": kitex.NewMethodInfo(listVideosHandler, newListVideosArgs, newListVideosResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "douyin.feed",
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

func listVideosHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(feed.ListFeedRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(feed.FeedService).ListVideos(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *ListVideosArgs:
		success, err := handler.(feed.FeedService).ListVideos(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ListVideosResult)
		realResult.Success = success
	}
	return nil
}
func newListVideosArgs() interface{} {
	return &ListVideosArgs{}
}

func newListVideosResult() interface{} {
	return &ListVideosResult{}
}

type ListVideosArgs struct {
	Req *feed.ListFeedRequest
}

func (p *ListVideosArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(feed.ListFeedRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ListVideosArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ListVideosArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ListVideosArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in ListVideosArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *ListVideosArgs) Unmarshal(in []byte) error {
	msg := new(feed.ListFeedRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ListVideosArgs_Req_DEFAULT *feed.ListFeedRequest

func (p *ListVideosArgs) GetReq() *feed.ListFeedRequest {
	if !p.IsSetReq() {
		return ListVideosArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ListVideosArgs) IsSetReq() bool {
	return p.Req != nil
}

type ListVideosResult struct {
	Success *feed.ListFeedResponse
}

var ListVideosResult_Success_DEFAULT *feed.ListFeedResponse

func (p *ListVideosResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(feed.ListFeedResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ListVideosResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ListVideosResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ListVideosResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in ListVideosResult")
	}
	return proto.Marshal(p.Success)
}

func (p *ListVideosResult) Unmarshal(in []byte) error {
	msg := new(feed.ListFeedResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ListVideosResult) GetSuccess() *feed.ListFeedResponse {
	if !p.IsSetSuccess() {
		return ListVideosResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ListVideosResult) SetSuccess(x interface{}) {
	p.Success = x.(*feed.ListFeedResponse)
}

func (p *ListVideosResult) IsSetSuccess() bool {
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

func (p *kClient) ListVideos(ctx context.Context, Req *feed.ListFeedRequest) (r *feed.ListFeedResponse, err error) {
	var _args ListVideosArgs
	_args.Req = Req
	var _result ListVideosResult
	if err = p.c.Call(ctx, "ListVideos", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}