package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/wechat"
	"toktik/logging"
	"toktik/service/wechat/db"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.EnvConfig.REDIS_ADDR,
		Password: config.EnvConfig.REDIS_PASSWORD,
		DB:       config.EnvConfig.REDIS_DB,
	})
}

// WechatServiceImpl implements the last service interface defined in the IDL.
type WechatServiceImpl struct{}

func (s *WechatServiceImpl) generateKey(sender, receiver *uint32) *string {
	key := fmt.Sprintf("chat:%d:%d", *sender, *receiver)
	return &key
}

// WechatChat implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatChat(ctx context.Context, req *wechat.MessageChatRequest) (resp *wechat.MessageChatResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"time":   time.Now(),
		"method": "WechatChat",
	})
	logger.Debugf("Process start")
	if req == nil {
		logger.Warningf("request is nil")
		return &wechat.MessageChatResponse{
			StatusCode: biz.RequestIsNil,
			StatusMsg:  "request is nil",
		}, nil
	}
	senderID := req.SenderId
	receiverID := req.ReceiverId
	key := *s.generateKey(&senderID, &receiverID)
	lRangeCMD := rdb.LRange(ctx, key, 0, -1)
	if lRangeCMD.Err() != nil {
		logger.Warningf("redis lrange error: %v", lRangeCMD.Err())
		resp = &wechat.MessageChatResponse{
			StatusCode: biz.RedisError,
			StatusMsg:  lRangeCMD.Err().Error(),
		}
		return
	}
	delCMD := rdb.Del(ctx, key)
	if delCMD.Err() != nil {
		logger.Warningf("redis del error: %v", delCMD.Err())
		resp = &wechat.MessageChatResponse{
			StatusCode: biz.RedisError,
			StatusMsg:  delCMD.Err().Error(),
		}
		return
	}
	respMessageList := make([]*wechat.Message, 0)
	messages := lRangeCMD.Val()
	for i, message := range messages {
		msg := &db.ChatMessage{}
		err = proto.Unmarshal([]byte(message), msg)
		if err != nil {
			logger.Warningf("proto unmarshal error: %v", err)
			resp = &wechat.MessageChatResponse{
				StatusCode: biz.ProtoUnmarshalError,
				StatusMsg:  err.Error(),
			}
			respMessageList = append(respMessageList, &wechat.Message{
				Id:         uint32(i),
				Content:    "[消息解析Message broken]",
				CreateTime: time.Unix(0, time.Now().UnixNano()).Format("2006-01-02 15:04:05"),
			})
		} else {
			respMessageList = append(respMessageList, &wechat.Message{
				Id:         uint32(i),
				Content:    msg.Msg,
				CreateTime: time.Unix(0, msg.Time*int64(time.Millisecond)).Format("2006-01-02 15:04:05"),
			})
		}
	}
	resp = &wechat.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: respMessageList,
	}
	return
}

// WechatAction implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatAction(ctx context.Context, req *wechat.MessageActionRequest) (resp *wechat.MessageActionResponse, err error) {
	startTime := time.Now()
	logger := logging.Logger.WithFields(logrus.Fields{
		"time":   startTime,
		"method": "WechatAction",
	})
	logger.Debugf("Process start")
	if req == nil {
		logger.Warningf("request is nil")
		return &wechat.MessageActionResponse{
			StatusCode: biz.RequestIsNil,
			StatusMsg:  "request is nil",
		}, nil
	}
	senderID := req.SenderId
	receiverID := req.ReceiverId
	msg := &db.ChatMessage{
		Msg:  req.Content,
		Time: startTime.UnixMilli(),
	}
	// get proto binary data
	msgStr, err := proto.Marshal(msg)
	if err != nil {
		logger.Warningf("proto marshal error: %v", err)
		resp = &wechat.MessageActionResponse{
			StatusCode: biz.ProtoMarshalError,
			StatusMsg:  err.Error(),
		}
		return
	}
	key := *s.generateKey(&senderID, &receiverID)
	cmd := rdb.LPush(ctx, key, msgStr)
	if cmd.Err() != nil {
		logger.Warningf("redis error: %v", cmd.Err())
		resp = &wechat.MessageActionResponse{
			StatusCode: biz.RedisError,
			StatusMsg:  cmd.Err().Error(),
		}
		return
	}
	return &wechat.MessageActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}, nil
}
