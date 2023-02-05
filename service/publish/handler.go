package main

import (
	"context"

	publish "toktik/kitex_gen/douyin/publish"
)

// PublishServiceImpl implements the last service interface defined in the IDL.
type PublishServiceImpl struct{}

// CreateVideo implements the PublishServiceImpl interface.
func (s *PublishServiceImpl) CreateVideo(ctx context.Context, req *publish.CreateVideoRequest) (resp *publish.CreateVideoResponse, err error) {
	return
}
