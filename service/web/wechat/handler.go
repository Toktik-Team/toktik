package wechat

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
	bizConstant "toktik/constant/biz"
	bizConfig "toktik/constant/config"
	"toktik/kitex_gen/douyin/wechat"
	"toktik/kitex_gen/douyin/wechat/wechatservice"
	"toktik/logging"
	"toktik/service/web/mw"
)

var Client wechatservice.Client

func init() {
	r, err := consul.NewConsulResolver(bizConfig.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	Client, err = wechatservice.NewClient(bizConfig.AuthServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func MessageAction(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"time":   time.Now(),
		"method": "MessageAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	userId := mw.GetAuthActorId(c)
	if userId == 0 {
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
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
		SenderId:   userId,
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
		"time":   time.Now(),
		"method": "MessageChat",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	userId := mw.GetAuthActorId(c)
	if userId == 0 {
		bizConstant.UnAuthorized.WithFields(&methodFields).LaunchError(c)
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

	logger.WithFields(logrus.Fields{
		"to_user_id": receiverIDInt,
	}).Debugf("Executing message chat")

	messageActionResponse, err := Client.WechatChat(ctx, &wechat.MessageChatRequest{
		SenderId:   userId,
		ReceiverId: uint32(receiverIDInt),
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
