package main

import (
	"context"
	"toktik/constant/config"
	"toktik/service/web/auth"
	"toktik/service/web/comment"
	"toktik/service/web/feed"
	"toktik/service/web/mw"
	"toktik/service/web/publish"
	"toktik/service/web/relation"
	"toktik/service/web/user"
	"toktik/service/web/wechat"

	"github.com/cloudwego/hertz/pkg/app"
	httpStatus "github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

func main() {
	h := server.Default(
		server.WithHostPorts(config.WebServiceAddr),
		server.WithMaxRequestBodySize(config.EnvConfig.MAX_REQUEST_BODY_SIZE),
	)
	h.Use(gzip.Gzip(gzip.DefaultCompression))
	h.Use(mw.AuthMiddleware())
	pprof.Register(h)

	douyin := h.Group("/douyin")

	// feed service
	douyin.GET("/feed", feed.Action)

	// user service
	userGroup := douyin.Group("/user")
	userGroup.POST("/register/", auth.Register)
	userGroup.POST("/login/", auth.Login)
	userGroup.GET("/", user.GetUserInfo)

	// publish service
	publishGroup := douyin.Group("/publish")
	publishGroup.POST("/action/", publish.Action)
	publishGroup.GET("/list", publish.List)

	// favorite service
	favoriteGroup := douyin.Group("/favorite")
	favoriteGroup.POST("/action", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(httpStatus.StatusOK, map[string]any{
			"status_code": 0,
			"message":     "ok",
		})
	})
	favoriteGroup.GET("/list", func(c context.Context, ctx *app.RequestContext) {
		ctx.JSON(httpStatus.StatusOK, map[string]any{
			"status_code": 0,
			"message":     "ok",
			"video_list":  []string{},
		})
	})

	// comment service
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action/", comment.Action)
	commentGroup.GET("/list/", comment.List)

	chatGPTUserList := map[string]any{
		"status_code": 0,
		"message":     "ok",
		"user_list": []map[string]any{{
			"id":             0,
			"name":           "ChatGPT",
			"follow_count":   1000000,
			"follower_count": 0,
			"is_follow":      true,
			"avatar":         "https://bkimg.cdn.bcebos.com/pic/8b13632762d0f703918f0d436fac463d269758ee6faf?x-bce-process=image/watermark,image_d2F0ZXIvYmFpa2U4MA==,g_7,xp_5,yp_5",
		}},
	}
	// relation service
	relationGroup := douyin.Group("/relation")
	relationGroup.POST("/action/", relation.RelationAction)
	relationGroup.GET("/follow/list/", relation.GetFollowList)
	relationGroup.GET("/follower/list/", relation.GetFollowerList)
	relationGroup.GET("/friend/list/", relation.GetFriendList)

	// message service
	messageGroup := douyin.Group("/message")
	messageGroup.POST("/action/", wechat.MessageAction)
	messageGroup.GET("/chat/", wechat.MessageChat)

	url := swagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	h.Spin()
}
