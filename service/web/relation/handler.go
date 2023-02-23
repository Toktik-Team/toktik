package relation

import (
	"context"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"strconv"
	bizConstant "toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/relation"
	relationService "toktik/kitex_gen/douyin/relation/relationservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var relationClient relationService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.RelationServiceName),
		provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	relationClient, err = relationService.NewClient(
		config.RelationServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
}

// RelationAction [POST] /relation/action/
func RelationAction(ctx context.Context, c *app.RequestContext) {
	var actionTypeInt int
	var relationResp *relation.RelationActionResponse

	methodFields := logrus.Fields{
		"method": "RelationAction",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Debugf("Process start")

	var actorId uint32
	switch c.GetString(mw.AuthResultKey) {
	case mw.AUTH_RESULT_SUCCESS:
		actorId = c.GetUint32(mw.UserIdKey)
	default:
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
		return
	}

	actionType, exist := c.GetQuery("action_type")
	if !exist {
		bizConstant.InvalidActionType.WithFields(&methodFields).LaunchError(c)
		return
	}
	actionTypeInt, err := strconv.Atoi(actionType)
	if err != nil || (actionTypeInt != 1 && actionTypeInt != 2) {
		err2Launch := bizConstant.InvalidActionType.WithFields(&methodFields)
		if err != nil {
			err2Launch = err2Launch.WithCause(err)
		}
		err2Launch.LaunchError(c)
		return
	}

	userId, exist := c.GetQuery("to_user_id")
	if !exist {
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
		return
	}
	userIdInt, err := strconv.ParseInt(userId, 10, 32)
	if err != nil {
		bizConstant.InvalidArguments.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	switch actionTypeInt {
	case 1:
		logger.WithFields(logrus.Fields{
			"actorId": actorId,
			"userId":  userId,
		}).Debugf("Executing follow")
		relationResp, err = relationClient.Follow(ctx, &relation.RelationActionRequest{
			UserId:  uint32(userIdInt),
			ActorId: actorId,
		})
	case 2:
		logger.WithFields(logrus.Fields{
			"userId":  userId,
			"actorId": actorId,
		}).Debugf("Executing unfollow")
		relationResp, err = relationClient.Unfollow(ctx, &relation.RelationActionRequest{
			UserId:  uint32(userIdInt),
			ActorId: actorId,
		})
	default:
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
		return
	}

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": relationResp,
	}).Debugf("Relation action success")
	c.JSON(
		httpStatus.StatusOK,
		relationResp,
	)
}

// GetFollowList FollowList [POST] /relation/follow/list/
func GetFollowList(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "GetFollowList",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Info("Process start")

	var actorId uint32
	switch c.GetString(mw.AuthResultKey) {
	case mw.AUTH_RESULT_SUCCESS, mw.AUTH_RESULT_NO_TOKEN:
		actorId = c.GetUint32(mw.UserIdKey)
	default:
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
		return
	}

	userId, idExist := c.GetQuery("user_id")
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		bizConstant.InvalidUserID.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithFields(logrus.Fields{
		"actorId": actorId,
		"userId":  id,
	}).Debugf("Executing get follow list")
	followListResp, err := relationClient.GetFollowList(ctx, &relation.FollowListRequest{
		UserId:  uint32(id),
		ActorId: actorId,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": followListResp,
	}).Debugf("get follow list success")
	c.JSON(
		httpStatus.StatusOK,
		followListResp,
	)
}

// GetFollowerList FollowList [POST] /relation/follower/list/
func GetFollowerList(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "GetFollowerList",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Info("Process start")

	var actorId uint32
	switch c.GetString(mw.AuthResultKey) {
	case mw.AUTH_RESULT_SUCCESS, mw.AUTH_RESULT_NO_TOKEN:
		actorId = c.GetUint32(mw.UserIdKey)
	default:
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
		return
	}

	userId, idExist := c.GetQuery("user_id")
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		bizConstant.InvalidUserID.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithFields(logrus.Fields{
		"actorId": actorId,
		"userId":  id,
	}).Debugf("Executing get follower list")
	followerListResp, err := relationClient.GetFollowerList(ctx, &relation.FollowerListRequest{
		UserId:  uint32(id),
		ActorId: actorId,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": followerListResp,
	}).Debugf("get follower list success")
	c.JSON(
		httpStatus.StatusOK,
		followerListResp,
	)
}

// GetFriendList [POST] /relation/friends/list/
func GetFriendList(ctx context.Context, c *app.RequestContext) {

	methodFields := logrus.Fields{
		"method": "GetFriendList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	var actorId uint32
	switch c.GetString(mw.AuthResultKey) {
	case mw.AUTH_RESULT_SUCCESS, mw.AUTH_RESULT_NO_TOKEN:
		actorId = c.GetUint32(mw.UserIdKey)
	default:
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
		return
	}

	userId, idExist := c.GetQuery("user_id")
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		bizConstant.InvalidUserID.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithFields(logrus.Fields{
		"userId":  userId,
		"actorId": id,
	}).Debugf("Executing get friend list")
	friendListResp, err := relationClient.GetFriendList(ctx, &relation.FriendListRequest{
		UserId:  uint32(id),
		ActorId: actorId,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": friendListResp,
	}).Debugf("get friends list success")
	c.JSON(
		httpStatus.StatusOK,
		friendListResp,
	)
}
