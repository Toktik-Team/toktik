package main

import (
	"toktik/constant/config"
	"toktik/service/web/auth"
	"toktik/service/web/comment"
	"toktik/service/web/feed"
	"toktik/service/web/mw"
	"toktik/service/web/publish"
	"toktik/service/web/user"

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

	douyin.Any("/authenticate", auth.Authenticate)

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
	publishGroup.GET("/list")

	// favorite service
	favoriteGroup := douyin.Group("/favorite")
	favoriteGroup.POST("/action")
	favoriteGroup.GET("/list")

	// comment service
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action/", comment.Action)
	commentGroup.GET("/list/", comment.List)

	// relation service
	relationGroup := douyin.Group("/relation")
	relationGroup.GET("/follow/list")
	relationGroup.GET("/follower/list")
	relationGroup.GET("/friend/list")

	// message service
	messageGroup := douyin.Group("/message")
	messageGroup.POST("/action")
	messageGroup.GET("/chat")

	url := swagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	h.Spin()
}
