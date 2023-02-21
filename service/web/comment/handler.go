package comment

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/comment/commentservice"
	"toktik/logging"
)

var commentClient commentservice.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	commentClient, err = commentservice.NewClient(config.CommentServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func Action(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"time":   time.Now(),
		"method": "CommentAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	actorId := c.GetUint32("user_id")
	videoId, videoIdExists := c.GetQuery("video_id")
	actionType, actionTypeExists := c.GetQuery("action_type")
	commentText, commentTextExists := c.GetQuery("comment_text")
	commentId, commentIdExists := c.GetQuery("comment_id")

	if actorId == 0 {
		biz.UnauthorizedError.WithFields(&methodFields).LaunchError(c)
		return
	}

	if !videoIdExists || !actionTypeExists {
		biz.BadRequestError.WithFields(&methodFields).LaunchError(c)
		return
	}

	pVideoId, err := strconv.ParseUint(videoId, 10, 32)
	pActionType, err := strconv.ParseInt(actionType, 10, 32)

	if err != nil {
		biz.BadRequestError.WithFields(&methodFields).WithCause(err).LaunchError(c)
		return
	}

	var rActionType = comment.ActionCommentType(pActionType)

	var (
		rErr error
	)

	switch rActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		if !commentTextExists {
			rErr = errors.New("comment text is required")
			break
		}
		resp, err := commentClient.ActionComment(ctx, &comment.ActionCommentRequest{
			ActorId:    actorId,
			VideoId:    uint32(pVideoId),
			ActionType: rActionType,
			Action:     &comment.ActionCommentRequest_CommentText{CommentText: commentText},
		})
		if err != nil {
			rErr = err
			break
		}
		logger.WithFields(logrus.Fields{
			"response": resp,
		}).Debugf("add comment success")
		c.JSON(
			consts.StatusOK,
			resp,
		)
		return
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		if !commentIdExists {
			rErr = errors.New("comment id is required")
			break
		}
		pCommentId, err := strconv.ParseUint(commentId, 10, 32)
		if err != nil {
			rErr = err
			break
		}
		resp, err := commentClient.ActionComment(ctx, &comment.ActionCommentRequest{
			ActorId:    actorId,
			VideoId:    uint32(pVideoId),
			ActionType: rActionType,
			Action:     &comment.ActionCommentRequest_CommentId{CommentId: uint32(pCommentId)},
		})
		if err != nil {
			rErr = err
			break
		}
		logger.WithFields(logrus.Fields{
			"response": resp,
		}).Debugf("delete comment success")
		c.JSON(
			consts.StatusOK,
			resp,
		)
		return
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		rErr = errors.New(fmt.Sprintf("invalid action type: %d", rActionType))
	}

	if rErr != nil {
		biz.InternalServerError.WithCause(rErr).WithFields(&methodFields).LaunchError(c)
		return
	}
}

func List(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"time":   time.Now(),
		"method": "CommentList",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	actorId := c.GetUint32("user_id")
	videoId, videoIdExists := c.GetQuery("video_id")

	if !videoIdExists {
		biz.BadRequestError.WithFields(&methodFields).LaunchError(c)
		return
	}

	pVideoId, err := strconv.ParseUint(videoId, 10, 32)

	if err != nil {
		biz.BadRequestError.WithFields(&methodFields).WithCause(err).LaunchError(c)
		return
	}

	resp, err := commentClient.ListComment(ctx, &comment.ListCommentRequest{
		ActorId: actorId,
		VideoId: uint32(pVideoId),
	})

	if err != nil {
		biz.InternalServerError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	c.JSON(
		consts.StatusOK,
		resp,
	)
}
