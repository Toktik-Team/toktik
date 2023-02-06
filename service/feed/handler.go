package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	gen "toktik/repo"
	"toktik/storage"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

const (
	OkStatusCode = 0
)

var (
	OkStatusMsg = "OK"
)

// ListVideos implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) ListVideos(ctx context.Context, req *feed.ListFeedRequest) (resp *feed.ListFeedResponse, err error) {
	publish := gen.Q.Video

	latestTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
	if err != nil {
		latestTime = time.Now().UnixMilli()
	}

	find, err := publish.WithContext(ctx).Where(publish.CreatedAt.Lte(time.UnixMilli(latestTime))).Order(publish.CreatedAt.Desc()).Limit(30).Offset(0).Find()
	if err != nil {
		return nil, err
	}

	nextTime := find[len(find)].CreatedAt.UnixMilli()

	var videos []*feed.Video
	for _, m := range find {

		u := &user.User{
			Id: uint32(m.UserId),
			// TODO: fill other fields
		}

		playUrl, err := storage.GetLink(m.FileName)
		if err != nil {
			_ = fmt.Errorf("failed to fetch play url: %w", err)
			continue
		}

		coverUrl, err := storage.GetLink(m.CoverName)
		if err != nil {
			_ = fmt.Errorf("failed to fetch cover url: %w", err)
			continue
		}

		videos = append(videos, &feed.Video{
			Id:       m.ID,
			Author:   u,
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
		StatusCode: OkStatusCode,
		StatusMsg:  &OkStatusMsg,
		NextTime:   &nextTime,
		Videos:     videos,
	}, nil
}
