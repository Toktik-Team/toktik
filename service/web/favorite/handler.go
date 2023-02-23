package favorite

import (
	"context"
	"net/http"
	"strconv"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/favorite"
	favoriteService "toktik/kitex_gen/douyin/favorite/favoriteservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var Client favoriteService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.FavoriteServiceName),
		provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	Client, err = favoriteService.NewClient(
		config.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
}

var logger = logging.Logger

// 用于解析 Action 函数所需参数
func parseParameters(c *app.RequestContext) (videoId uint32, actionType uint32, isEnd bool) {
	field := logrus.Fields{
		"method": "parseParameters",
	}

	isEnd = true
	// 获取参数
	qVideoId, videoIdExist := c.GetQuery("video_id")
	qActionType, actionTypeExist := c.GetQuery("action_type")
	if !videoIdExist || !actionTypeExist {
		biz.BadRequestError.
			WithFields(&field).
			LaunchError(c)
		return
	}
	// 解析 videoId
	temp, err := strconv.ParseUint(qVideoId, 10, 32)
	if err != nil {
		biz.BadRequestError.
			WithCause(err).
			WithFields(&field).
			LaunchError(c)
		return
	}
	videoId = uint32(temp)
	// 解析 actionType
	temp, err = strconv.ParseUint(qActionType, 10, 32)
	if err != nil {
		biz.BadRequestError.
			WithCause(err).
			WithFields(&field).
			LaunchError(c)
		return
	}
	actionType = uint32(temp)
	if actionType != 1 && actionType != 2 {
		biz.BadRequestError.
			WithFields(&field).
			LaunchError(c)
		return
	}

	isEnd = false
	return
}

// Action 处理点赞和取消点赞
func Action(ctx context.Context, c *app.RequestContext) {
	field := logrus.Fields{
		"method": "Action",
	}
	logger.WithFields(field).Debugf("Process start")

	actorIdPtr, ok := mw.Auth(c, mw.WithAuthRequired())
	actorId := *actorIdPtr
	if !ok {
		return
	}

	videoId, actionType, isEnd := parseParameters(c)
	if isEnd {
		return
	}

	response, err := Client.FavoriteAction(ctx, &favorite.FavoriteRequest{
		ActorId:    actorId,
		VideoId:    videoId,
		ActionType: actionType,
	})

	if err != nil {
		biz.RPCCallError.WithCause(err).WithFields(&field).LaunchError(c)
		return
	}

	c.JSON(
		http.StatusOK,
		response,
	)
	return
}

// List 列出用户所有点赞视频
func List(ctx context.Context, c *app.RequestContext) {
	field := logrus.Fields{
		"method": "List",
	}
	logger.WithFields(field).Info("Process start")

	actorIdPtr, ok := mw.Auth(c)
	actorId := *actorIdPtr
	if !ok {
		return
	}

	qUserId, userIdExist := c.GetQuery("user_id")
	if !userIdExist {
		biz.BadRequestError.WithFields(&field).LaunchError(c)
		return
	}

	userId, err := strconv.ParseUint(qUserId, 10, 32)
	if err != nil {
		biz.BadRequestError.WithCause(err).WithFields(&field).LaunchError(c)
		return
	}

	response, err := Client.FavoriteList(ctx, &favorite.FavoriteListRequest{
		ActorId: actorId,
		UserId:  uint32(userId),
	})

	if err != nil {
		biz.RPCCallError.WithCause(err).WithFields(&field).LaunchError(c)
		return
	}

	c.JSON(
		consts.StatusOK,
		response,
	)
	return
}
