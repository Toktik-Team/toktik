package user

import (
	"context"
	"strconv"
	bizConstant "toktik/constant/biz"
	bizConfig "toktik/constant/config"
	"toktik/kitex_gen/douyin/user"
	userService "toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var userClient userService.Client

func init() {
	r, err := consul.NewConsulResolver(bizConfig.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	userClient, err = userService.NewClient(bizConfig.UserServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "GetUserInfo",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	userId, idExist := c.GetQuery("user_id")
	actorId := mw.GetAuthActorId(c)
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		bizConstant.InvalidUserID.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithField("user_id", id).Debugf("Executing get user info")
	resp, err := userClient.GetUser(ctx, &user.UserRequest{
		UserId:  uint32(id),
		ActorId: actorId,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithField("response", resp).Debugf("Get user info success")
	c.JSON(
		httpStatus.StatusOK,
		resp,
	)
}
