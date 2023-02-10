package main

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	gen "toktik/repo"
	"toktik/storage"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// ListVideos implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) ListVideos(ctx context.Context, req *feed.ListFeedRequest) (resp *feed.ListFeedResponse, err error) {
	video := gen.Q.Video

	latestTime, err := strconv.ParseInt(*req.LatestTime, 10, 64)
	if err != nil {
		if _, ok := err.(*strconv.NumError); ok {
			latestTime = time.Now().UnixMilli()
		} else {
			resp = &feed.ListFeedResponse{
				StatusCode: biz.Unable2ParseLatestTimeStatusCode,
				StatusMsg:  &biz.BadRequestStatusMsg,
				NextTime:   nil,
				Videos:     nil,
			}
			return resp, nil
		}
	}

	find, err := video.WithContext(ctx).
		Where(video.CreatedAt.Lte(time.UnixMilli(latestTime))).
		Order(video.CreatedAt.Desc()).
		Limit(biz.VideoCount).
		Offset(0).
		Find()
	if err != nil {
		// TODO: handle error
		return nil, err
	}

	nextTime := find[len(find)].CreatedAt.UnixMilli()

	var videos []*feed.Video
	for _, m := range find {

		u := &user.User{
			Id: m.UserId,
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
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		NextTime:   &nextTime,
		Videos:     videos,
	}, nil
}