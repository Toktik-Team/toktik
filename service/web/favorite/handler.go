package favorite

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"toktik/constant/biz"
	bizConfig "toktik/constant/config"
	"toktik/kitex_gen/douyin/favorite"
	favoriteService "toktik/kitex_gen/douyin/favorite/favoriteservice"
	"toktik/logging"
	"toktik/service/web/mw"
)

var Client favoriteService.Client

func init() {
	r, err := consul.NewConsulResolver(bizConfig.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	Client, err = favoriteService.NewClient(bizConfig.FavoriteServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

var logger = logging.Logger

// 获取 actorId
func getActorId(c *app.RequestContext, actorId *uint32) bool {
	authResult := mw.GetAuthResult(c)
	switch authResult {
	case mw.AUTH_RESULT_SUCCESS:
		*actorId = mw.GetAuthActorId(c)
		return false
	default:
		biz.AuthFailed.
			WithFields(&logrus.Fields{
				"method": "GetActorId",
			}).
			LaunchError(c)
		return true
	}
}

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

	var actorId uint32
	if getActorId(c, &actorId) {
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

	if mw.GetAuthResult(c) != mw.AUTH_RESULT_SUCCESS {
		biz.AuthFailed.WithFields(&field).LaunchError(c)
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
		UserId: uint32(userId),
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
