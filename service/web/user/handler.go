package user

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"strconv"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/user"
	"toktik/kitex_gen/douyin/user/userservice"
)

var userClient userservice.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	userClient, err = userservice.NewClient(config.UserServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	userId, idExist := c.GetQuery("user_id")
	requesterId := c.GetUint32("user_id")
	id, err := strconv.Atoi(userId)

	if !idExist || err != nil {
		c.JSON(
			consts.StatusBadRequest,
			&user.UserResponse{
				StatusCode: 1,
				StatusMsg:  "invalid request",
				User:       nil,
			},
		)
		return
	}

	resp, err := userClient.GetUser(ctx, &user.UserRequest{
		UserId:      uint32(id),
		RequesterId: requesterId,
	})

	if err != nil {
		c.JSON(
			consts.StatusOK,
			&user.UserResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
				User:       nil,
			},
		)
		return
	}

	c.JSON(
		consts.StatusOK,
		resp,
	)
}
