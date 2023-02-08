package main

import (
	"context"
	user "toktik/kitex_gen/douyin/user"
	"toktik/repo"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	userToken := repo.UserToken
	_, err = userToken.WithContext(ctx).Where(userToken.Token.Eq(req.Token)).First()

	if err != nil {
		resp = &user.UserResponse{
			StatusCode: 1,
			StatusMsg:  "user not logged in",
			User:       nil,
		}
		return
	}

	userInfo := repo.User
	u, err := userInfo.WithContext(ctx).Where(userInfo.ID.Eq(req.UserId)).First()

	if err != nil {
		resp = &user.UserResponse{
			StatusCode: 1,
			StatusMsg:  "user does exist",
			User:       nil,
		}
		return
	}

	resp.User = &user.User{
		Id:            u.ID,
		Name:          u.Username,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      false, //TODO: 是否关注
	}
	return
}
