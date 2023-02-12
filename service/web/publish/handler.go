package publish

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"toktik/constant/config"
	"toktik/kitex_gen/douyin/publish"
	publishService "toktik/kitex_gen/douyin/publish/publishservice"
	"toktik/service/web/mw"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
)

var publishClient publishService.Client

func init() {
	r, err := consul.NewConsulResolver(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	publishClient, err = publishService.NewClient(config.PublishServiceName, client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
	}
}

func paramValidate(ctx context.Context, c *app.RequestContext) (err error) {
	form, err := c.Request.MultipartForm()
	if err != nil {
		return fmt.Errorf("invalid form: %w", err)
	}
	title := form.Value["title"]
	if len(title) <= 0 {
		return fmt.Errorf("not title")
	}

	data := form.File["data"]
	if len(data) <= 0 {
		return fmt.Errorf("not data")
	}
	return nil
}

func Action(ctx context.Context, c *app.RequestContext) {
	if err := paramValidate(ctx, c); err != nil {
		c.JSON(
			consts.StatusBadRequest,
			&publish.CreateVideoResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		)
		return
	}

	form, _ := c.MultipartForm()
	title := form.Value["title"][0]
	file := form.File["data"]
	opened, _ := file[0].Open()
	defer func(opened multipart.File) {
		err := opened.Close()
		if err != nil {
			log.Println(err)
		}
	}(opened)
	var data = make([]byte, file[0].Size)
	readSize, err := opened.Read(data)
	if err != nil {
		c.JSON(
			consts.StatusBadRequest,
			&publish.CreateVideoResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		)
		return
	}
	if readSize != int(file[0].Size) {
		c.JSON(
			consts.StatusBadRequest,
			&publish.CreateVideoResponse{
				StatusCode: 1,
				StatusMsg:  "read size not equal to file size",
			},
		)
		return
	}
	userId := c.GetUint32(mw.USER_ID_KEY)
	publishResp, err := publishClient.CreateVideo(ctx, &publish.CreateVideoRequest{
		UserId: userId,
		Data:   data,
		Title:  title,
	})
	if err != nil {
		c.JSON(
			consts.StatusOK,
			&publish.CreateVideoResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		)
		return
	}
	c.JSON(
		consts.StatusOK,
		publishResp,
	)
}
