package comment

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"strconv"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/comment/commentservice"
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
	actorId := c.GetUint32("user_id")
	videoId, videoIdExists := c.GetQuery("video_id")
	actionType, actionTypeExists := c.GetQuery("action_type")
	commentText := c.Query("comment_text")
	commentId := c.Query("comment_id")

	if actorId == 0 {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	if !videoIdExists || !actionTypeExists {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	pVideoId, err := strconv.ParseUint(videoId, 10, 32)
	pActionType, err := strconv.ParseInt(actionType, 10, 32)
	pCommentId, err := strconv.ParseUint(commentId, 10, 32)

	if err != nil {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	var rActionType = comment.ActionCommentType(pActionType)

	var (
		rErr error
	)

	switch rActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
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
		c.JSON(
			consts.StatusOK,
			resp,
		)
		return
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
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
		c.JSON(
			consts.StatusInternalServerError,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.InternalServerErrorStatusMsg,
			},
		)
		return
	}
}

func List(ctx context.Context, c *app.RequestContext) {
	actorId := c.GetUint32("user_id")
	videoId, videoIdExists := c.GetQuery("video_id")

	if actorId == 0 {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	if !videoIdExists {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	pVideoId, err := strconv.ParseUint(videoId, 10, 32)

	if err != nil {
		c.JSON(
			consts.StatusBadRequest,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.BadRequestStatusMsg,
			},
		)
		return
	}

	resp, err := commentClient.ListComment(ctx, &comment.ListCommentRequest{
		ActorId: actorId,
		VideoId: uint32(pVideoId),
	})

	if err != nil {
		c.JSON(
			consts.StatusInternalServerError,
			&comment.ActionCommentResponse{
				StatusCode: 1,
				StatusMsg:  &biz.InternalServerErrorStatusMsg,
			},
		)
		return
	}

	c.JSON(
		consts.StatusOK,
		resp,
	)
}
