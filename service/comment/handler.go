package main

import (
	"context"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
	"log"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	gen "toktik/repo"
	"toktik/repo/model"
)

var UserClient userservice.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	UserClient, err = userservice.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// ActionComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ActionComment(ctx context.Context, req *comment.ActionCommentRequest) (resp *comment.ActionCommentResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"user_id":      req.ActorId,
		"video_id":     req.VideoId,
		"action_type":  req.ActionType,
		"comment_text": req.GetCommentText(),
		"comment_id":   req.GetCommentId(),
		"time":         time.Now(),
		"function":     "ActionComment",
	})
	logger.Debug("Process start")

	var pCommentText string
	var pCommentID uint32
	switch req.ActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		pCommentText = req.GetCommentText()
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		pCommentID = req.GetCommentId()
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_UNSPECIFIED:
		fallthrough
	default:
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("invalid action type")
		return &comment.ActionCommentResponse{
			StatusCode: biz.InvalidCommentActionTypeStatusCode,
			StatusMsg:  &biz.BadRequestStatusMsg,
		}, nil
	}

	// Video check: check if the qVideo exists || check if creator is the same as actor
	qVideo := gen.Q.Video
	pVideo, err := qVideo.WithContext(ctx).
		Where(qVideo.ID.Eq(req.VideoId)).
		First()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("video query error")
		return &comment.ActionCommentResponse{
			StatusCode: biz.UnableToQueryVideoStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}
	if pVideo == nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("video not found")
		return &comment.ActionCommentResponse{
			StatusCode: biz.VideoNotFoundStatusCode,
			StatusMsg:  &biz.BadRequestStatusMsg,
		}, nil
	}

	userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
		UserId:  req.ActorId,
		ActorId: req.ActorId,
	})

	if err != nil || userResponse.StatusCode != biz.OkStatusCode {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("user service error")
		return &comment.ActionCommentResponse{
			StatusCode: biz.InternalUserServiceErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	pUser := userResponse.User

	switch req.ActionType {
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_ADD:
		resp, err = addComment(ctx, logger, pUser, pVideo.ID, pCommentText)
	case comment.ActionCommentType_ACTION_COMMENT_TYPE_DELETE:
		resp, err = deleteComment(ctx, logger, pUser, pVideo.ID, pCommentID)
	}
	if err != nil {
		return resp, err
	}

	logger.WithFields(logrus.Fields{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")

	return resp, err
}

func addComment(ctx context.Context, logger *logrus.Entry, pUser *user.User, pVideoID uint32, pCommentText string) (resp *comment.ActionCommentResponse, err error) {
	count, err := gen.Q.Comment.WithContext(ctx).Count()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed to query db entry")
		resp = &comment.ActionCommentResponse{
			StatusCode: biz.UnableToQueryCommentStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}

	rComment := model.Comment{
		VideoId:   pVideoID,
		CommentId: uint32(count + 1),
		UserId:    pUser.Id,
		Content:   pCommentText,
	}

	err = gen.Q.Comment.WithContext(ctx).Create(&rComment)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed to create db entry")
		resp = &comment.ActionCommentResponse{
			StatusCode: biz.UnableToCreateCommentStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}
	resp = &comment.ActionCommentResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Comment: &comment.Comment{
			Id:         rComment.CommentId,
			User:       pUser,
			Content:    rComment.Content,
			CreateDate: rComment.CreatedAt.Format("01-02"),
		},
	}
	return
}

func deleteComment(ctx context.Context, logger *logrus.Entry, pUser *user.User, videoID uint32, commentID uint32) (resp *comment.ActionCommentResponse, err error) {
	qComment := gen.Q.Comment

	rComment, err := qComment.WithContext(ctx).
		Where(qComment.VideoId.Eq(videoID), qComment.CommentId.Eq(commentID)).
		First()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed to query db entry")
		resp = &comment.ActionCommentResponse{
			StatusCode: biz.UnableToQueryCommentStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}

	if rComment.UserId != pUser.Id {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("comment creator and actor not match")
		return &comment.ActionCommentResponse{
			StatusCode: biz.ActorIDNotMatchStatusCode,
			StatusMsg:  &biz.ForbiddenStatusMsg,
		}, nil
	}

	_, err = qComment.WithContext(ctx).Delete(rComment)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
			"err":  err,
		}).Debug("failed to delete db entry")
		resp = &comment.ActionCommentResponse{
			StatusCode: biz.UnableToDeleteCommentStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}
	resp = &comment.ActionCommentResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Comment:    nil,
	}
	return
}

// ListComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ListComment(ctx context.Context, req *comment.ListCommentRequest) (resp *comment.ListCommentResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"user_id":  req.ActorId,
		"video_id": req.VideoId,
		"time":     time.Now(),
		"function": "ListComment",
	})
	logger.Debug("Process start")

	qVideo := gen.Q.Video
	pVideo, err := qVideo.WithContext(ctx).
		Where(qVideo.ID.Eq(req.VideoId)).
		First()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("video query error")
		return &comment.ListCommentResponse{
			StatusCode: biz.UnableToQueryVideoStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}
	if pVideo == nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("video not found")
		return &comment.ListCommentResponse{
			StatusCode: biz.VideoNotFoundStatusCode,
			StatusMsg:  &biz.BadRequestStatusMsg,
		}, nil
	}

	qComment := gen.Q.Comment
	pCommentList, err := qComment.WithContext(ctx).
		Where(qComment.VideoId.Eq(pVideo.ID)).
		Order(qComment.CreatedAt.Desc()).
		Find()

	if err != nil {
		logger.WithFields(logrus.Fields{
			"time": time.Now(),
		}).Debug("comment query error")
		return &comment.ListCommentResponse{
			StatusCode: biz.UnableToQueryCommentStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	rCommentList := make([]*comment.Comment, len(pCommentList))
	for _, pComment := range pCommentList {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  req.ActorId,
			ActorId: req.ActorId,
		})

		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			logger.WithFields(logrus.Fields{
				"pComment": pComment,
				"err":      err,
				"time":     time.Now(),
			}).Debug("unable to get user info")
		}

		rCommentList = append(rCommentList, &comment.Comment{
			Id:         pComment.CommentId,
			User:       userResponse.User,
			Content:    pComment.Content,
			CreateDate: pComment.CreatedAt.Format("01-02"),
		})
	}

	logger.WithFields(logrus.Fields{
		"time":     time.Now(),
		"response": resp,
	}).Debug("all process done, ready to launch response")
	return &comment.ListCommentResponse{
		StatusCode:  biz.OkStatusCode,
		StatusMsg:   &biz.OkStatusMsg,
		CommentList: rCommentList,
	}, nil
}
