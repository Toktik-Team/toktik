package auth

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"strconv"
	"toktik/config"
	"toktik/kitex_gen/douyin/auth"
	authService "toktik/kitex_gen/douyin/auth/authservice"
	"toktik/kitex_gen/douyin/feed"
	feedService "toktik/kitex_gen/douyin/feed/feedservice"
	"toktik/kitex_gen/douyin/publish"
	publishService "toktik/kitex_gen/douyin/publish/publishservice"
)

var AuthClient authService.Client
var publishClient publishService.Client
var feedClient feedService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	AuthClient, err = authService.NewClient(config.AuthServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
	publishClient, err = publishService.NewClient(config.PublishServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
	feedClient, err = feedService.NewClient(config.FeedServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func Authenticate(ctx context.Context, c *app.RequestContext) {
	token, exist := c.GetQuery("token")
	if !exist {
		c.String(consts.StatusUnauthorized, "failed")
		return
	}
	authenticateResp, err := AuthClient.Authenticate(ctx, &auth.AuthenticateRequest{Token: token})
	if err != nil {
		c.String(consts.StatusUnauthorized, "failed")
		return
	}
	c.String(consts.StatusOK, strconv.Itoa(int(authenticateResp.UserId)))
}

func PublishAction(ctx context.Context, c *app.RequestContext) {
	// TODO: read user id from gateway
	userId := 1
	publishResp, err := publishClient.CreateVideo(ctx, &publish.CreateVideoRequest{
		UserId: uint32(userId),
	})
	if err != nil {
		c.JSON(
			consts.StatusOK,
			&publish.CreateVideoResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		)
		return
	}
	c.JSON(
		consts.StatusOK,
		publishResp,
	)
	return
}

func FeedAction(ctx context.Context, c *app.RequestContext) {
	latestTime := c.Query("latest_time")

	token := c.Query("token")

	response, err := feedClient.ListVideos(ctx, &feed.ListFeedRequest{
		LatestTime: &latestTime,
		Token:      &token,
	})
	if err != nil {
		c.JSON(
			consts.StatusOK,
			struct {
				StatusCode    int    `json:"status_code"`
				StatusMessage string `json:"status_message"`
			}{1, err.Error()},
		)
		return
	}
	c.JSON(
		consts.StatusOK,
		response,
	)
}
