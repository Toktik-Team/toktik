package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/comment"
	commentService "toktik/kitex_gen/douyin/comment/commentservice"
	"toktik/kitex_gen/douyin/favorite"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	userService "toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	gen "toktik/repo"
	"toktik/repo/model"
	"toktik/storage"
)

var UserClient userService.Client
var CommentClient commentService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	UserClient, err = userService.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
	CommentClient, err = commentService.NewClient(config.CommentServiceName, client.WithResolver(r))
	if err != nil {
		panic(err)
	}
}

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

var logger = logging.Logger

func like(ctx context.Context, actorId uint32, videoId uint32, authorId uint32) (resp *favorite.FavoriteResponse, err error) {
	q := gen.Use(gen.DB)
	err = q.Transaction(func(tx *gen.Query) error {
		// 加入用户喜爱列表
		err := tx.Favorite.WithContext(ctx).Create(&model.Favorite{
			UserId:  actorId,
			VideoId: videoId,
		})
		if err != nil {
			return err
		}
		// 增加用户总点赞数
		_, err = tx.User.WithContext(ctx).
			Where(tx.User.ID.Eq(actorId)).
			Update(tx.User.FavoriteCount, tx.User.FavoriteCount.Add(1))
		if err != nil {
			return err
		}
		// 增加作者总获赞数
		_, err = tx.User.WithContext(ctx).
			Where(tx.User.ID.Eq(authorId)).
			Update(tx.User.FavoriteCount, tx.User.FavoriteCount.Add(1))
		if err != nil {
			return err
		}
		// 增加视频总获赞数
		_, err = tx.Video.WithContext(ctx).
			Where(tx.Video.ID.Eq(videoId)).
			Update(tx.Video.FavoriteCount, tx.Video.FavoriteCount.Add(1))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.UnableToLike,
			StatusMsg:  nil,
		}

		logger.WithFields(logrus.Fields{
			"method": "like",
		}).Errorf("Failed to like: %v", err)

		return
	}

	resp = &favorite.FavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
	}
	return
}

func cancelLike(ctx context.Context, actorId uint32, videoId uint32, authorId uint32) (resp *favorite.FavoriteResponse, err error) {
	q := gen.Use(gen.DB)
	err = q.Transaction(func(tx *gen.Query) error {
		// 从用户喜爱列表中移除
		_, err := tx.Favorite.WithContext(ctx).Delete(&model.Favorite{
			UserId:  actorId,
			VideoId: videoId,
		})
		if err != nil {
			return err
		}
		// 减少用户总点赞数
		_, err = tx.User.WithContext(ctx).
			Where(tx.User.ID.Eq(actorId)).
			Update(tx.User.FavoriteCount, tx.User.FavoriteCount.Sub(1))
		if err != nil {
			return err
		}
		// 减少作者总获赞数
		_, err = tx.User.WithContext(ctx).
			Where(tx.User.ID.Eq(authorId)).
			Update(tx.User.FavoriteCount, tx.User.FavoriteCount.Sub(1))
		if err != nil {
			return err
		}
		// 减少视频总获赞数
		_, err = tx.Video.WithContext(ctx).
			Where(tx.Video.ID.Eq(videoId)).
			Update(tx.Video.FavoriteCount, tx.Video.FavoriteCount.Sub(1))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.UnableToCancelLike,
			StatusMsg:  nil,
		}

		logger.WithFields(logrus.Fields{
			"method": "cancelLike",
		}).Errorf("Failed to cancel like: %v", err)

		return
	}

	resp = &favorite.FavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
	}
	return
}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteRequest) (resp *favorite.FavoriteResponse, err error) {
	field := logrus.Fields{
		"method": "FavoriteAction",
	}
	logger.WithFields(field).Info("process start")

	v := gen.Video
	var authorId uint32
	err = v.WithContext(ctx).Where(v.ID.Eq(req.VideoId)).Pluck(v.UserId, &authorId)
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.UnableToQueryVideo,
			StatusMsg:  nil,
		}

		logger.WithFields(field).Error("Failed to get video author")
		return
	}

	if req.ActionType == 1 {
		resp, err = like(ctx, req.ActorId, req.VideoId, authorId)
	} else {
		resp, err = cancelLike(ctx, req.ActorId, req.VideoId, authorId)
	}

	return
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	field := logrus.Fields{
		"method": "FavoriteList",
	}
	logger.WithFields(field).Info("process start")

	f := gen.Favorite
	var videoIds []uint32
	err = f.WithContext(ctx).Where(f.UserId.Eq(req.UserId)).Pluck(f.VideoId, &videoIds)
	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: biz.FailedToGetVideoList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			VideoList:  nil,
		}
		logger.WithFields(field).
			Errorf("Failed to get the id of the user's favorite videos: %v", err)
		return
	}

	v := gen.Video
	videos, err := v.WithContext(ctx).Where(v.ID.In(videoIds...)).Find()
	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: biz.FailedToGetVideoList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			VideoList:  nil,
		}
		logger.WithFields(field).
			Errorf("Failed to get information of user's favorite videos: %v", err)
		return
	}

	var videoList []*feed.Video

	for _, video := range videos {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  video.UserId,
			ActorId: req.UserId,
		})
		if err != nil {
			logger.
				WithFields(field).
				WithFields(logrus.Fields{
					"function called": "UserClient.GetUser",
					"parameters":      fmt.Sprintf("UserId: %d, ActorId: %d", video.UserId, req.UserId),
				}).Errorf("Failed to get user information: %v", err)
			continue
		}

		playUrl, err := storage.GetLink(video.FileName)
		if err != nil {
			logger.
				WithFields(field).
				WithFields(logrus.Fields{
					"function called": "storage.GetLink",
					"parameters":      fmt.Sprintf("fileName: %s", video.FileName),
				}).Errorf("Failed to get play url: %v", err)
			continue
		}

		coverUrl, err := storage.GetLink(video.CoverName)
		if err != nil {
			logger.
				WithFields(field).
				WithFields(logrus.Fields{
					"function called": "storage.GetLink",
					"parameters":      fmt.Sprintf("fileName: %s", video.CoverName),
				}).Errorf("Failed to get play url: %v", err)
			continue
		}

		commentCount, err := CommentClient.CountComment(ctx, &comment.CountCommentRequest{
			ActorId: req.UserId,
			VideoId: video.ID,
		})

		if err != nil {
			logger.
				WithFields(field).
				WithFields(logrus.Fields{
					"function called": "CommentClient.CountComment",
					"parameters":      fmt.Sprintf("ActorId: %d, VideoId: %d", req.UserId, video.ID),
				}).Errorf("Failed to get the number of comments: %v", err)
			continue
		}

		videoList = append(videoList, &feed.Video{
			Id:            video.ID,
			Author:        userResponse.User,
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  commentCount.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}

	resp = &favorite.FavoriteListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  videoList,
	}
	return
}

// IsFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) IsFavorite(ctx context.Context, req *favorite.IsFavoriteRequest) (resp *favorite.IsFavoriteResponse, err error) {
	field := logrus.Fields{
		"method": "IsFavorite",
	}
	logger.WithFields(field).Info("process start")

	f := gen.Favorite
	_, err = f.WithContext(ctx).Where(f.UserId.Eq(req.UserId), f.VideoId.Eq(req.VideoId)).First()
	if err != nil {
		logger.
			WithFields(field).
			WithFields(logrus.Fields{
				"parameters": fmt.Sprintf("UserId: %v, VideoId: %d", req.UserId, req.VideoId),
			}).Error(err.Error())
		return &favorite.IsFavoriteResponse{Result: false}, err
	}

	return &favorite.IsFavoriteResponse{Result: true}, nil
}
