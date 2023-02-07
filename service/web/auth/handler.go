package auth

import (
	"context"
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
		c.String(consts.StatusUnauthorized, "failed")
		return
	}
	authenticateResp, err := AuthClient.Authenticate(ctx, &auth.AuthenticateRequest{Token: token})
	if err != nil {
		c.String(consts.StatusUnauthorized, "failed")
		return
	}
	c.String(consts.StatusOK, strconv.Itoa(int(authenticateResp.UserId)))
}
