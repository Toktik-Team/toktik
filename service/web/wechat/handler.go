package wechat

import (
	"context"
	"log"
	"strconv"
	bizConstant "toktik/constant/biz"
	bizConfig "toktik/constant/config"
	"toktik/kitex_gen/douyin/wechat"
	"toktik/kitex_gen/douyin/wechat/wechatservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var Client wechatservice.Client

func init() {
	r, err := consul.NewConsulResolver(bizConfig.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(bizConfig.WechatServiceName),
		provider.WithExportEndpoint(bizConfig.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	Client, err = wechatservice.NewClient(
		bizConfig.WechatServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func MessageAction(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "MessageAction",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Debugf("Process start")

	actorIdPtr, ok := mw.Auth(c, mw.WithAuthRequired())
	actorId := *actorIdPtr
	if !ok {
		return
	}
	actionType, exist := c.GetQuery("action_type")
	if !exist {
		bizConstant.InvalidActionType.WithFields(&methodFields).LaunchError(c)
		return
	}
	if i, err := strconv.Atoi(actionType); err != nil || i != 1 {
		err2Launch := bizConstant.InvalidActionType.WithFields(&methodFields)
		if err != nil {
			err2Launch = err2Launch.WithCause(err)
		}
		err2Launch.LaunchError(c)
		return
	}

	receiverID, exist := c.GetQuery("to_user_id")
	if !exist {
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
		return
	}
	receiverIDInt, err := strconv.ParseInt(receiverID, 10, 32)
	if err != nil {
		bizConstant.InvalidArguments.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	content, exist := c.GetQuery("content")
	if !exist {
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithFields(logrus.Fields{
		"action_type": actionType,
		"to_user_id":  receiverIDInt,
		"content":     content,
	}).Debugf("Executing message action")

	messageActionResponse, err := Client.WechatAction(ctx, &wechat.MessageActionRequest{
		SenderId:   actorId,
		ReceiverId: uint32(receiverIDInt),
		ActionType: 1,
		Content:    content,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": messageActionResponse,
	}).Debugf("Message action success")
	c.JSON(httpStatus.StatusOK, messageActionResponse)
}

func MessageChat(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "MessageChat",
	}
	logger := logging.Logger
	logger.WithFields(methodFields).Debugf("Process start")

	actorIdPtr, ok := mw.Auth(c, mw.WithAuthRequired())
	actorId := *actorIdPtr
	if !ok {
		return
	}

	receiverIdPtr, exist := c.GetQuery("to_user_id")
	if !exist {
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
		return
	}
	receiverId, err := strconv.ParseInt(receiverIdPtr, 10, 32)
	if err != nil {
		bizConstant.InvalidArguments.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	preMsgTimeStr, exist := c.GetQuery("pre_msg_time")
	if !exist {
		preMsgTimeStr = "0"
	}
	preMsgTime, err := strconv.ParseInt(preMsgTimeStr, 10, 64)
	if err != nil {
		bizConstant.InvalidArguments.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	logger.WithFields(logrus.Fields{
		"to_user_id":   receiverId,
		"pre_msg_time": preMsgTimeStr,
	}).Debugf("Executing message chat")

	messageActionResponse, err := Client.WechatChat(ctx, &wechat.MessageChatRequest{
		SenderId:   actorId,
		ReceiverId: uint32(receiverId),
		PreMsgTime: preMsgTime,
	})

	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": messageActionResponse,
	}).Debugf("Message chat success")
	c.JSON(httpStatus.StatusOK, messageActionResponse)
}
