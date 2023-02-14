package main

import (
	"context"
	"toktik/kitex_gen/douyin/user"
	"toktik/repo"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// GetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	if req == nil {
		resp = &user.UserResponse{
			StatusCode: 0,
			StatusMsg:  "user does exist",
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
			StatusCode: 1,
			StatusMsg:  "user does exist",
			User:       nil,
		}
		return
	}
	resp = &user.UserResponse{}

	resp.User = &user.User{
		Id:            u.ID,
		Name:          u.Username,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      false, //TODO: 是否关注
	}
	return
}
