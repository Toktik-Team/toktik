package main

import (
	"context"
	"google.golang.org/protobuf/proto"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/favorite"
	"toktik/kitex_gen/douyin/feed"
	feedService "toktik/kitex_gen/douyin/feed/feedservice"
	"toktik/logging"
	gen "toktik/repo"
	"toktik/repo/model"

	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var FeedClient feedService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	FeedClient, err = feedService.NewClient(
		config.FeedServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
}

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

var logger = logging.Logger

func like(ctx context.Context, actorId uint32, videoId uint32) (resp *favorite.FavoriteResponse, err error) {
	if err := gen.Q.Favorite.WithContext(ctx).Create(&model.Favorite{
		UserId:  actorId,
		VideoId: videoId,
	}); err != nil {

		resp = &favorite.FavoriteResponse{
			StatusCode: biz.UnableToLike,
			StatusMsg:  proto.String("点赞失败"),
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

func cancelLike(ctx context.Context, actorId uint32, videoId uint32) (resp *favorite.FavoriteResponse, err error) {
	if _, err = gen.Q.Favorite.WithContext(ctx).
		Where(gen.Favorite.UserId.Eq(actorId), gen.Favorite.VideoId.Eq(videoId)).
		Unscoped().
		Delete(); err != nil {

		resp = &favorite.FavoriteResponse{
			StatusCode: biz.UnableToCancelLike,
			StatusMsg:  proto.String("取消点赞失败"),
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
			StatusMsg:  proto.String("视频作者不存在"),
		}

		logger.WithFields(field).Error("Failed to get video author")
		return
	}

	if req.ActionType == 1 {
		resp, err = like(ctx, req.ActorId, req.VideoId)
	} else {
		resp, err = cancelLike(ctx, req.ActorId, req.VideoId)
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

	queryVideoResp, err := FeedClient.QueryVideos(ctx, &feed.QueryVideosRequest{
		ActorId:  req.ActorId,
		VideoIds: videoIds,
	})
	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: biz.FailedToGetVideoList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			VideoList:  nil,
		}
		logger.WithFields(field).
			Errorf("Failed to get the user's favorite videos: %v", err)
		return
	}

	resp = &favorite.FavoriteListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  queryVideoResp.VideoList,
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
		return &favorite.IsFavoriteResponse{Result: false}, nil
	}

	return &favorite.IsFavoriteResponse{Result: true}, nil
}

// CountFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) CountFavorite(ctx context.Context, req *favorite.CountFavoriteRequest) (resp *favorite.CountFavoriteResponse, err error) {
	field := logrus.Fields{
		"method": "CountFavorite",
	}
	logger.WithFields(field).Info("process start")

	count, err := gen.Q.Favorite.WithContext(ctx).
		Where(gen.Q.Favorite.VideoId.Eq(req.VideoId)).
		Count()
	if err != nil {
		logger.WithFields(field).
			Errorf("Failed to get the number of favorites: %v", err)
		return &favorite.CountFavoriteResponse{
			StatusCode: biz.SQLQueryErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	return &favorite.CountFavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}

// CountUserFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) CountUserFavorite(ctx context.Context, req *favorite.CountUserFavoriteRequest) (resp *favorite.CountUserFavoriteResponse, err error) {
	field := logrus.Fields{
		"method": "CountUserFavorite",
	}
	logger.WithFields(field).Info("process start")

	count, err := gen.Q.Favorite.WithContext(ctx).
		Where(gen.Q.Favorite.UserId.Eq(req.UserId)).
		Count()
	if err != nil {
		logger.WithFields(field).
			Errorf("Failed to get the number of favorites: %v", err)
		return &favorite.CountUserFavoriteResponse{
			StatusCode: biz.SQLQueryErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	return &favorite.CountUserFavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}

// CountUserTotalFavorited implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) CountUserTotalFavorited(ctx context.Context, req *favorite.CountUserTotalFavoritedRequest) (resp *favorite.CountUserTotalFavoritedResponse, err error) {
	field := logrus.Fields{
		"method": "CountUserTotalFavorited",
	}
	logger.WithFields(field).Info("process start")

	videos, err := gen.Q.Video.WithContext(ctx).
		Where(gen.Q.Video.UserId.Eq(req.UserId)).
		Find()
	if err != nil {
		logger.WithFields(field).
			Errorf("Failed to get the number of favorites: %v", err)
		return &favorite.CountUserTotalFavoritedResponse{
			StatusCode: biz.SQLQueryErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	videoIds := make([]uint32, 0, len(videos))
	for _, video := range videos {
		videoIds = append(videoIds, video.ID)
	}

	count, err := gen.Q.Favorite.WithContext(ctx).
		Where(gen.Q.Favorite.VideoId.In(videoIds...)).
		Count()
	if err != nil {
		logger.WithFields(field).
			Errorf("Failed to get the number of favorites: %v", err)
		return &favorite.CountUserTotalFavoritedResponse{
			StatusCode: biz.SQLQueryErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}, nil
	}

	return &favorite.CountUserTotalFavoritedResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		Count:      uint32(count),
	}, nil
}
