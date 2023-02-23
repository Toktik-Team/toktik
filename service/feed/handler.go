package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/comment"
	"toktik/kitex_gen/douyin/comment/commentservice"
	"toktik/kitex_gen/douyin/favorite"
	favoriteService "toktik/kitex_gen/douyin/favorite/favoriteservice"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
	"toktik/logging"
	gen "toktik/repo"
	"toktik/repo/model"
	"toktik/storage"

	"github.com/cloudwego/kitex/client"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var UserClient userservice.Client
var CommentClient commentservice.Client
var FavoriteClient favoriteService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	UserClient, err = userservice.NewClient(
		config.UserServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		log.Fatal(err)
	}
	CommentClient, err = commentservice.NewClient(
		config.CommentServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		log.Fatal(err)
	}
	FavoriteClient, err = favoriteService.NewClient(
		config.FavoriteServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// ListVideos implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) ListVideos(ctx context.Context, req *feed.ListFeedRequest) (resp *feed.ListFeedResponse, err error) {
	methodFields := logrus.Fields{
		"latest_time": req.LatestTime,
		"actor_id":    req.ActorId,
		"function":    "ListVideos",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	now := time.Now().UnixMilli()
	latestTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			latestTime = now
		} else {
			resp = &feed.ListFeedResponse{
				StatusCode: biz.Unable2ParseLatestTimeStatusCode,
				StatusMsg:  &biz.BadRequestStatusMsg,
				NextTime:   nil,
				VideoList:  nil,
			}
			return resp, nil
		}
	}

	find, err := findVideos(ctx, latestTime)
	if err != nil {
		resp = &feed.ListFeedResponse{
			StatusCode: biz.SQLQueryErrorStatusCode,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			NextTime:   &now,
			VideoList:  nil,
		}
		return resp, nil
	}

	if len(find) == 0 {
		resp = &feed.ListFeedResponse{
			StatusCode: biz.OkStatusCode,
			StatusMsg:  &biz.OkStatusMsg,
			NextTime:   nil,
			VideoList:  nil,
		}
		return resp, nil
	}
	nextTime := find[len(find)-1].CreatedAt.Add(time.Duration(-1)).UnixMilli()

	var actorId uint32 = 0
	if req.ActorId != nil {
		actorId = *req.ActorId
	}
	videos := queryDetailed(ctx, logger, actorId, find)

	return &feed.ListFeedResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		NextTime:   &nextTime,
		VideoList:  videos,
	}, nil
}

func findVideos(ctx context.Context, latestTime int64) ([]*model.Video, error) {
	video := gen.Q.Video
	return video.WithContext(ctx).
		Where(video.CreatedAt.Lte(time.UnixMilli(latestTime))).
		Order(video.CreatedAt.Desc()).
		Limit(biz.VideoCount).
		Offset(0).
		Find()
}

// QueryVideos implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) QueryVideos(ctx context.Context, req *feed.QueryVideosRequest) (resp *feed.QueryVideosResponse, err error) {
	methodFields := logrus.Fields{
		"actor_id":  req.ActorId,
		"video_ids": req.VideoIds,
		"function":  "ListVideos",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debug("Process start")

	rst, err := query(ctx, logger, req.ActorId, req.VideoIds)
	if err != nil {
		resp = &feed.QueryVideosResponse{
			StatusCode: biz.UnableToQueryUser,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			VideoList:  rst,
		}
		return resp, nil
	}

	return &feed.QueryVideosResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  rst,
	}, nil
}

func queryDetailed(
	ctx context.Context,
	logger *logrus.Entry,
	actorId uint32,
	videos []*model.Video,
) (respVideoList []*feed.Video) {
	wg := sync.WaitGroup{}
	respVideoList = make([]*feed.Video, len(videos))
	for i, v := range videos {
		respVideoList[i] = &feed.Video{
			Id:     v.ID,
			Title:  v.Title,
			Author: &user.User{Id: v.UserId},
		}
		wg.Add(6)
		// fill author
		go func(i int, v *model.Video) {
			defer wg.Done()
			userResponse, localErr := UserClient.GetUser(ctx, &user.UserRequest{
				UserId:  v.UserId,
				ActorId: actorId,
			})
			if localErr != nil || userResponse.StatusCode != biz.OkStatusCode {
				logger.WithFields(logrus.Fields{
					"video_id": v.ID,
					"user_id":  v.UserId,
					"cause":    localErr,
				}).Warning("failed to get user info")
				return
			}
			respVideoList[i].Author = userResponse.User
		}(i, v)

		// fill play url
		go func(i int, v *model.Video) {
			defer wg.Done()
			playUrl, localErr := storage.GetLink(v.FileName)
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"video_id":  v.ID,
					"file_name": v.FileName,
					"err":       localErr,
				}).Warning("failed to fetch play url")
				return
			}
			respVideoList[i].PlayUrl = playUrl
		}(i, v)

		// fill cover url
		go func(i int, v *model.Video) {
			defer wg.Done()
			coverUrl, localErr := storage.GetLink(v.CoverName)
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"video_id":   v.ID,
					"cover_name": v.CoverName,
					"err":        localErr,
				}).Warning("failed to fetch cover url")
				return
			}
			respVideoList[i].CoverUrl = coverUrl
		}(i, v)

		// fill favorite count
		go func(i int, v *model.Video) {
			defer wg.Done()
			favoriteCount, localErr := FavoriteClient.CountFavorite(ctx, &favorite.CountFavoriteRequest{
				VideoId: v.ID,
			})
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"video_id": v.ID,
					"err":      localErr,
				}).Warning("failed to fetch favorite count")
				return
			}
			respVideoList[i].FavoriteCount = favoriteCount.Count
		}(i, v)

		// fill comment count
		go func(i int, v *model.Video) {
			defer wg.Done()
			commentCount, localErr := CommentClient.CountComment(ctx, &comment.CountCommentRequest{
				ActorId: actorId,
				VideoId: v.ID,
			})
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"video_id": v.ID,
					"err":      localErr,
				}).Warning("failed to fetch comment count")
				return
			}
			respVideoList[i].CommentCount = commentCount.CommentCount
		}(i, v)

		// fill is favorite
		go func(i int, v *model.Video) {
			defer wg.Done()
			isFavorite, localErr := FavoriteClient.IsFavorite(ctx, &favorite.IsFavoriteRequest{
				UserId:  actorId,
				VideoId: v.ID,
			})
			if localErr != nil {
				logger.WithFields(logrus.Fields{
					"video_id": v.ID,
					"err":      localErr,
				}).Warning("failed to fetch favorite status")
				return
			}
			respVideoList[i].IsFavorite = isFavorite.Result
		}(i, v)
	}
	wg.Wait()

	return
}

func query(ctx context.Context, logger *logrus.Entry, actorId uint32, videoIds []uint32) (resp []*feed.Video, err error) {
	find, err := gen.Q.Video.WithContext(ctx).Where(gen.Q.Video.ID.In(videoIds...)).Find()
	if err != nil {
		return nil, err
	}

	return queryDetailed(ctx, logger, actorId, find), nil
}
