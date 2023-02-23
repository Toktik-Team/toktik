package user

import (
	"context"
	"strconv"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/user"
	userService "toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var userClient userService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.UserServiceName),
		provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	userClient, err = userService.NewClient(
		config.UserServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "GetUserInfo",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Info("Process start")

	actorIdPtr, ok := mw.Auth(c)
	actorId := *actorIdPtr
	if !ok {
		return
	}

	userId, idExist := c.GetQuery("user_id")
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		biz.InvalidUserID.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithField("user_id", id).Debugf("Executing get user info")
	resp, err := userClient.GetUser(ctx, &user.UserRequest{
		UserId:  uint32(id),
		ActorId: actorId,
	})

	if err != nil {
		biz.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithField("response", resp).Debugf("Get user info success")
	c.JSON(
		httpStatus.StatusOK,
		resp,
	)
}
