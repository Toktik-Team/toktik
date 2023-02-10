package feed

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/feed"
	feedService "toktik/kitex_gen/douyin/feed/feedservice"
)

var feedClient feedService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	feedClient, err = feedService.NewClient(config.FeedServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}
func Action(ctx context.Context, c *app.RequestContext) {
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
