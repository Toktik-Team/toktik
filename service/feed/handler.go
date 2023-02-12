package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"strconv"
	"time"
	"toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
	gen "toktik/repo"
	"toktik/repo/model"
	"toktik/storage"
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

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// ListVideos implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) ListVideos(ctx context.Context, req *feed.ListFeedRequest) (resp *feed.ListFeedResponse, err error) {
	latestTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			latestTime = time.Now().UnixMilli()
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
			NextTime:   nil,
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

	var videos []*feed.Video
	var actorId uint32 = 0
	if req.ActorId != nil {
		actorId = *req.ActorId
	}

	for _, m := range find {

		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  m.UserId,
			ActorId: actorId,
		})
		if err != nil || userResponse.StatusCode != biz.OkStatusCode {
			log.Println(fmt.Errorf("failed to get user info: %w", err))
			continue
		}

		playUrl, err := storage.GetLink(m.FileName)
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch play url: %w", err))
			continue
		}

		coverUrl, err := storage.GetLink(m.CoverName)
		if err != nil {
			log.Println(fmt.Errorf("failed to fetch cover url: %w", err))
			continue
		}

		videos = append(videos, &feed.Video{
			Id:       m.ID,
			Author:   userResponse.User,
			PlayUrl:  playUrl,
			CoverUrl: coverUrl,
			// TODO: finish this
			FavoriteCount: 0,
			// TODO: finish this
			CommentCount: 0,
			// TODO: finish this
			IsFavorite: false,
			Title:      m.Title,
		})
	}

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
