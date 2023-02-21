package main

import (
	"toktik/constant/config"
	"toktik/service/web/auth"
	"toktik/service/web/comment"
	"toktik/service/web/favorite"
	"toktik/service/web/feed"
	"toktik/service/web/mw"
	"toktik/service/web/publish"
	"toktik/service/web/relation"
	"toktik/service/web/user"
	"toktik/service/web/wechat"

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
	favoriteGroup.POST("/action/", favorite.Action)
	favoriteGroup.GET("/list/", favorite.List)

	// comment service
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action/", comment.Action)
	commentGroup.GET("/list/", comment.List)

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
