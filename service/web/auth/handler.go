package auth

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"strconv"
	"toktik/config"
	"toktik/kitex_gen/douyin/auth"
	authService "toktik/kitex_gen/douyin/auth/authservice"
)

var AuthClient authService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	AuthClient, err = authService.NewClient(config.AuthServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func Authenticate(ctx context.Context, c *app.RequestContext) {
	token, exist := c.GetQuery("token")
	if !exist {
		c.JSON(consts.StatusUnauthorized, &auth.AuthenticateResponse{
			StatusCode: 1,
			StatusMsg:  "no token",
		})
		return
	}
	authenticateResp, err := AuthClient.Authenticate(ctx, &auth.AuthenticateRequest{Token: token})
	if err != nil {
		c.JSON(consts.StatusUnauthorized, authenticateResp)
		return
	}
	c.String(consts.StatusOK, strconv.Itoa(int(authenticateResp.UserId)))
}

func Register(ctx context.Context, c *app.RequestContext) {
	username, usernameExist := c.GetQuery("username")
	password, passwordExist := c.GetQuery("password")
	if !usernameExist || !passwordExist {
		c.JSON(consts.StatusUnauthorized, &auth.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "no username or password",
		})
		return
	}
	registerResponse, err := AuthClient.Register(ctx, &auth.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		c.JSON(consts.StatusUnauthorized, &auth.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("unknown error: %v", err),
		})
		return
	}
	c.JSON(consts.StatusOK, registerResponse)
}

func Login(ctx context.Context, c *app.RequestContext) {
	username, usernameExist := c.GetQuery("username")
	password, passwordExist := c.GetQuery("password")
	if !usernameExist || !passwordExist {
		c.JSON(consts.StatusUnauthorized, &auth.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  "no username or password",
		})
		return
	}
	loginResponse, err := AuthClient.Login(ctx, &auth.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		c.JSON(consts.StatusUnauthorized, &auth.RegisterResponse{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("unknown error: %v", err),
		})
		return
	}
	c.JSON(consts.StatusOK, loginResponse)
}
