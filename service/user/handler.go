package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/user"
	"toktik/logging"
	"toktik/repo"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {

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

	resp = &user.UserResponse{
		StatusCode: biz.OkStatusCode,
		StatusMsg:  &biz.OkStatusMsg,
		User: &user.User{
			Id:              u.ID,
			Name:            u.Name,
			FollowCount:     u.FollowCount,
			FollowerCount:   u.FollowerCount,
			IsFollow:        false, // TODO: 是否关注
			Avatar:          &avatar,
			BackgroundImage: &backgroundImage,
			Signature:       &u.Name,          // TODO:
			TotalFavorited:  u.TotalFavorited, // TODO：获赞时更新获赞数量
			WorkCount:       u.WorkCount,
			FavoriteCount:   u.FavoriteCount, // TODO：点赞时更新点赞数量
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
