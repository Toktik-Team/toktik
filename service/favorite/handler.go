package main

import (
	"context"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"toktik/constant/biz"
	"toktik/constant/config"
	favorite "toktik/kitex_gen/douyin/favorite"
	"toktik/kitex_gen/douyin/feed"
	"toktik/kitex_gen/douyin/user"
	userService "toktik/kitex_gen/douyin/user/userservice"
	gen "toktik/repo"
	"toktik/repo/model"
	"toktik/storage"
)

var UserClient userService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	UserClient, err = userService.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

func like(ctx context.Context, actorId uint32, videoId uint32) (resp *favorite.FavoriteResponse, err error) {
	u := gen.Q.User
	v := gen.Q.Video
	m := model.User{Model: model.Model{ID: actorId}}
	video := model.Video{Model: model.Model{ID: videoId}}
	// 加入用户喜爱列表
	err = u.FavoriteVideo.WithContext(ctx).Model(&m).Append(&video)
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.FailedToLikeVideo,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}
	// 增加视频总点赞数
	_, err = v.WithContext(ctx).Where(v.ID.Eq(videoId)).Update(v.FavoriteCount, v.FavoriteCount.Add(1))
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.FailedToAddVideoFavoriteCount,
			StatusMsg:  nil,
		}
		return
	}

	resp = &favorite.FavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
	}
	return
}

func cancelLike(ctx context.Context, actorId uint32, videoId uint32) (resp *favorite.FavoriteResponse, err error) {
	u := gen.Q.User
	v := gen.Q.Video
	m := model.User{Model: model.Model{ID: actorId}}
	video := model.Video{Model: model.Model{ID: videoId}}
	// 从用户喜爱列表中移除
	err = u.FavoriteVideo.WithContext(ctx).Model(&m).Delete(&video)
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.FailedToCancelLike,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
		return
	}
	// 减少视频总点赞数
	_, err = v.WithContext(ctx).Where(v.ID.Eq(videoId)).Update(v.FavoriteCount, v.FavoriteCount.Sub(1))
	if err != nil {
		resp = &favorite.FavoriteResponse{
			StatusCode: biz.FailedToSubVideoFavoriteCount,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
		}
	}

	resp = &favorite.FavoriteResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
	}
	return
}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteRequest) (resp *favorite.FavoriteResponse, err error) {
	if req.ActionType == 1 {
		resp, err = like(ctx, req.ActorId, req.VideoId)
	} else {
		resp, err = cancelLike(ctx, req.ActorId, req.VideoId)
	}

	return
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	u := gen.Q.User
	m := model.User{Model: model.Model{ID: req.UserId}}
	videos, err := u.FavoriteVideo.Model(&m).Find()

	if err != nil {
		resp = &favorite.FavoriteListResponse{
			StatusCode: biz.FailedToGetVideoList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			VideoList:  nil,
		}
		return
	}

	var videoList []*feed.Video

	for _, v := range videos {
		userResponse, err := UserClient.GetUser(ctx, &user.UserRequest{
			UserId:  v.UserId,
			ActorId: req.UserId,
		})
		if err != nil {

		}

		playUrl, err := storage.GetLink(v.FileName)
		if err != nil {

		}

		coverUrl, err := storage.GetLink(v.CoverName)
		if err != nil {

		}

		videoList = append(videoList, &feed.Video{
			Id:            v.ID,
			Author:        userResponse.User,
			PlayUrl:       playUrl,
			CoverUrl:      coverUrl,
			FavoriteCount: v.FavoriteCount,
			// TODO: 评论总数
			CommentCount: 0,
			IsFavorite:   true,
			Title:        v.Title,
		})
	}

	resp = &favorite.FavoriteListResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		VideoList:  videoList,
	}
	return
}
