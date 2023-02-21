package auth

import (
	"context"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	bizConstant "toktik/constant/biz"
	bizConfig "toktik/constant/config"
	"toktik/kitex_gen/douyin/auth"
	authService "toktik/kitex_gen/douyin/auth/authservice"
	"toktik/logging"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var Client authService.Client

func init() {
	r, err := consul.NewConsulResolver(bizConfig.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(bizConfig.AuthServiceName),
		provider.WithExportEndpoint(bizConfig.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	Client, err = authService.NewClient(
		bizConfig.AuthServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()))
	if err != nil {
		log.Fatal(err)
	}
}

func Register(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "Register",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	username, usernameExist := c.GetQuery("username")
	password, passwordExist := c.GetQuery("password")
	if !usernameExist || !passwordExist {
		bizConstant.NoUserNameOrPassWord.WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": username,
		"password": password,
	}).Debugf("Executing register")
	registerResponse, err := Client.Register(ctx, &auth.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": registerResponse,
	}).Debugf("Register success")
	c.JSON(httpStatus.StatusOK, registerResponse)
}

func Login(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "Login",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	username, usernameExist := c.GetQuery("username")
	password, passwordExist := c.GetQuery("password")
	if !usernameExist || !passwordExist {
		bizConstant.NoUserNameOrPassWord.WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"username": username,
		"password": password,
	}).Debugf("Executing login")
	loginResponse, err := Client.Login(ctx, &auth.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": loginResponse,
	}).Debugf("Login success")
	c.JSON(httpStatus.StatusOK, loginResponse)
}
