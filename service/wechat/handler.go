package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/wechat"
	"toktik/logging"
	"toktik/service/wechat/db"

	"github.com/golang/protobuf/proto"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protojson"
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

// WechatAction implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatAction(ctx context.Context, req *wechat.MessageActionRequest) (resp *wechat.MessageActionResponse, err error) {
	actionTime := time.Now()
	logger := logging.Logger.WithFields(logrus.Fields{
		"actionTime": actionTime,
		"method":     "WechatAction",
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

	// 0 means chat with GPT
	if receiverID == 0 {
		err := s.handleChatGPT(ctx, senderID, req.Content, false, actionTime.UnixMilli())
		if err != nil {
			logger.Warningf("handleChatGPT error: %v", err)
			return &wechat.MessageActionResponse{
				StatusCode: biz.ServiceNotAvailable,
				StatusMsg:  err.Error(),
			}, nil
		}
		return &wechat.MessageActionResponse{
			StatusCode: 0,
			StatusMsg:  "success",
		}, nil
	}

	msg := &db.ChatMessage{
		From: senderID,
		To:   receiverID,
		Msg:  req.Content,
		Time: actionTime.UnixMilli(),
	}

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
	cmd := rdb.ZAdd(ctx, key, redis.Z{
		Score:  float64(actionTime.UnixMilli()),
		Member: msgStr,
	})
	if cmd.Err() != nil {
		logger.Warningf("redis error: %v", cmd.Err())
		resp = &wechat.MessageActionResponse{
			StatusCode: biz.RedisError,
			StatusMsg:  cmd.Err().Error(),
		}
		return
	}
	resp = &wechat.MessageActionResponse{
		StatusCode: 0,
		StatusMsg:  "success",
	}
	logger.WithFields(logrus.Fields{
		"sender_id":   senderID,
		"receiver_id": receiverID,
		"content":     req.Content,
		"cost_time":   time.Since(actionTime).Milliseconds(),
	}).Debugf("Process end")
	return
}

// WechatChat implements the WechatServiceImpl interface.
func (s *WechatServiceImpl) WechatChat(ctx context.Context, req *wechat.MessageChatRequest) (resp *wechat.MessageChatResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
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
	min := strconv.FormatInt(int64(req.PreMsgTime), 10)
	max := "+inf"
	key := *s.generateKey(&senderID, &receiverID)
	zRangeCMD := rdb.ZRangeByScore(ctx, key, &redis.ZRangeBy{Min: fmt.Sprintf("%f", min), Max: max})
	if zRangeCMD.Err() != nil {
		logger.Warningf("redis lrange error: %v", zRangeCMD.Err())
		resp = &wechat.MessageChatResponse{
			StatusCode: biz.RedisError,
			StatusMsg:  zRangeCMD.Err().Error(),
		}
		return
	}
	respMessageList := make([]*wechat.Message, 0)
	messages := zRangeCMD.Val()
	for i, message := range messages {
		msg := &db.ChatMessage{}
		err = proto.Unmarshal([]byte(message), msg)
		var content string
		if err != nil {
			content = fmt.Sprintf("%s", message)
			logger.Warningf("proto unmarshal error: %v", err)
		} else {
			content = msg.Msg
		}
		respMessageList = append(respMessageList, &wechat.Message{
			Id:         uint32(i),
			Content:    content,
			CreateTime: strconv.FormatInt(msg.Time, 10),
			FromUserId: &msg.From,
			ToUserId:   &msg.To,
		})
	}
	resp = &wechat.MessageChatResponse{
		StatusCode:  0,
		StatusMsg:   "success",
		MessageList: respMessageList,
	}
	logger.WithFields(logrus.Fields{
		"message_list": respMessageList,
	}).Debugf("Process end")
	return resp, nil
}

func (s *WechatServiceImpl) handleChatGPT(ctx context.Context, senderID uint32, content string, resetSession bool, time int64) (err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"time":   time,
		"method": "handleChatGPT",
	})
	logger.Debugf("Process start")
	chatGPTMessage := db.ChatGPTMessage{
		SenderId:     senderID,
		Msg:          content,
		ResetSession: resetSession,
		Time:         time,
	}
	msgStr, err := protojson.Marshal(&chatGPTMessage)
	if err != nil {
		logger.Warningf("proto marshal error: %v", err)
		return err
	}
	cmd := rdb.Publish(ctx, "chatgpt", msgStr)
	if cmd.Err() != nil {
		logger.Warningf("redis publish error: %v", cmd.Err())
		return cmd.Err()
	}
	return
}
