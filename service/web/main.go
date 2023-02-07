package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"toktik/config"
	"toktik/service/web/auth"
	"toktik/service/web/mw"
)

func main() {
	h := server.Default(server.WithHostPorts(config.WebServiceAddr))
	h.Use(gzip.Gzip(gzip.DefaultCompression))
	h.Use(mw.ProtoJsonMiddleware())
	h.Use(mw.AuthMiddleware())
	pprof.Register(h)

	douyin := h.Group("/douyin")

	douyin.Any("/authenticate", auth.Authenticate)

	// feed service
	douyin.GET("/feed")

	// user service
	userGroup := douyin.Group("/user")
	userGroup.POST("/register", auth.Register)
	userGroup.POST("/login", auth.Login)
	userGroup.GET("/")

	// publish service
	publishGroup := douyin.Group("/publish")
	publishGroup.POST("/action")
	publishGroup.GET("/list")

	// favorite service
	favoriteGroup := douyin.Group("/favorite")
	favoriteGroup.POST("/action")
	favoriteGroup.GET("/list")

	// comment service
	commentGroup := douyin.Group("/comment")
	commentGroup.POST("/action")
	commentGroup.GET("/list")

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
