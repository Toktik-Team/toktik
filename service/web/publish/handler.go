package publish

import (
	"context"
	"fmt"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"mime/multipart"
	"strconv"
	bizConstant "toktik/constant/biz"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/publish"
	publishService "toktik/kitex_gen/douyin/publish/publishservice"
	"toktik/logging"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/sirupsen/logrus"
)

var publishClient publishService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.PublishServiceName),
		provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)
	publishClient, err = publishService.NewClient(
		config.PublishServiceName,
		client.WithResolver(r),
		client.WithSuite(tracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
}

func paramValidate(c *app.RequestContext) (err error) {
	var wrappedError error
	form, err := c.Request.MultipartForm()
	if err != nil {
		wrappedError = fmt.Errorf("invalid form: %w", err)
	}
	title := form.Value["title"]
	if len(title) <= 0 {
		wrappedError = fmt.Errorf("not title")
	}

	data := form.File["data"]
	if len(data) <= 0 {
		wrappedError = fmt.Errorf("not data")
	}
	if wrappedError != nil {
		return wrappedError
	}
	return nil
}

func Action(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "PublishAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	if err := paramValidate(c); err != nil {
		bizConstant.InvalidArguments.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	form, _ := c.MultipartForm()
	title := form.Value["title"][0]
	file := form.File["data"][0]
	opened, _ := file.Open()
	defer func(opened multipart.File) {
		err := opened.Close()
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("opened.Close() failed")
		}
	}(opened)
	var data = make([]byte, file.Size)
	readSize, err := opened.Read(data)
	if err != nil {
		bizConstant.OpenFileFailedError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	if readSize != int(file.Size) {
		bizConstant.SizeNotMatchError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	userId := mw.GetAuthActorId(c)

	logger.WithFields(logrus.Fields{
		"userId":    userId,
		"title":     title,
		"data_size": len(data),
	}).Debugf("Executing create video")
	publishResp, err := publishClient.CreateVideo(ctx, &publish.CreateVideoRequest{
		UserId: userId,
		Data:   data,
		Title:  title,
	})
	if err != nil {
		bizConstant.RPCCallError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}
	logger.WithFields(logrus.Fields{
		"response": publishResp,
	}).Debugf("Create video success")
	c.JSON(
		httpStatus.StatusOK,
		publishResp,
	)
}

func List(ctx context.Context, c *app.RequestContext) {
	methodFields := logrus.Fields{
		"method": "CommentAction",
	}
	logger := logging.Logger.WithFields(methodFields)
	logger.Debugf("Process start")

	actorId := c.GetUint32("user_id")
	userId, userIdExists := c.GetQuery("user_id")

	if actorId == 0 {
		bizConstant.UnauthorizedError.WithFields(&methodFields).LaunchError(c)
		return
	}

	if !userIdExists {
		bizConstant.InvalidArguments.WithFields(&methodFields).LaunchError(c)
	}

	pUserId, err := strconv.ParseUint(userId, 10, 32)

	if err != nil {
		bizConstant.BadRequestError.WithFields(&methodFields).WithCause(err).LaunchError(c)
		return
	}

	resp, err := publishClient.ListVideo(ctx, &publish.ListVideoRequest{
		UserId:  uint32(pUserId),
		ActorId: actorId,
	})

	if err != nil {
		bizConstant.InternalServerError.WithCause(err).WithFields(&methodFields).LaunchError(c)
		return
	}

	c.JSON(
		httpStatus.StatusOK,
		resp,
	)
}
