package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/favorite"
	favoriteService "toktik/kitex_gen/douyin/favorite/favoriteservice"
	"toktik/kitex_gen/douyin/publish"
	publishService "toktik/kitex_gen/douyin/publish/publishservice"
	"toktik/kitex_gen/douyin/relation"
	relationService "toktik/kitex_gen/douyin/relation/relationservice"
	"toktik/kitex_gen/douyin/user"
	"toktik/logging"
	"toktik/repo"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

var FavoriteClient favoriteService.Client
var RelationClient relationService.Client
var PublishClient publishService.Client

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO: 等到 kitex 更新后删除此代码
	if req == nil {
		resp = &user.UserResponse{
			StatusCode: biz.UserNotFound,
			User: &user.User{
				Id:            0,
				Name:          "anonymous",
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
			},
		}
		return
	}

	userInfo := repo.User
	u, err := userInfo.WithContext(ctx).Where(userInfo.ID.Eq(req.UserId)).First()

	if err != nil {
		resp = &user.UserResponse{
			StatusCode: biz.UserNotFound,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	avatar := u.GetUserAvatar()
	backgroundImage := u.GetBackgroundImage()

	followCount, err := RelationClient.CountFollowList(ctx, &relation.CountFollowListRequest{
		UserId: u.ID,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user follow count failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryFollowList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	followerCount, err := RelationClient.CountFollowerList(ctx, &relation.CountFollowerListRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user follower count failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryFollowerList,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	isFollow, err := RelationClient.IsFollow(ctx, &relation.IsFollowRequest{
		ActorId: req.ActorId,
		UserId:  req.UserId,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user is follow failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryIsFollow,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	totalFavorited, err := FavoriteClient.UserTotalFavoritedCount(ctx, &favorite.UserTotalFavoritedCountRequest{
		ActorId: req.ActorId,
		UserId:  req.UserId,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user total favorited failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryTotalFavorited,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	workCount, err := PublishClient.CountVideo(ctx, &publish.CountVideoRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user work count failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryVideo,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	favoriteCount, err := FavoriteClient.UserFavoriteCount(ctx, &favorite.UserFavoriteCountRequest{
		UserId: req.UserId,
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("get user favorite count failed")
		resp = &user.UserResponse{
			StatusCode: biz.UnableToQueryFavorite,
			StatusMsg:  &biz.InternalServerErrorStatusMsg,
			User:       nil,
		}
		return
	}

	resp = &user.UserResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		User: &user.User{
			Id:              u.ID,
			Name:            u.Username,
			FollowCount:     followCount.Count,
			FollowerCount:   followerCount.Count,
			IsFollow:        isFollow.Result,
			Avatar:          &avatar,
			BackgroundImage: &backgroundImage,
			Signature:       &u.Username, // TODO
			TotalFavorited:  &totalFavorited.Count,
			WorkCount:       &workCount.Count,
			FavoriteCount:   &favoriteCount.Count,
		},
	}
	if u.IsUpdated() {
		_, err = userInfo.WithContext(ctx).Where(userInfo.ID.Eq(u.ID)).Updates(u)
		if err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"time": time.Now(),
				"err":  err,
			}).Errorf("save user failed")
		}
	}
	return
}
