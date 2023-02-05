package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/pprof"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
	"toktik/config"
)

func main() {
	h := server.Default(server.WithHostPorts(config.WebServiceAddr))
	h.Use(gzip.Gzip(gzip.DefaultCompression))
	h.Use(ProtoJsonMiddleware())
	pprof.Register(h)

	h.Any("/authenticate", Authenticate)

	// feed service
	h.GET("/feed")

	// user service
	userGroup := h.Group("/user")
	userGroup.POST("/register")
	userGroup.POST("/login")
	userGroup.GET("/")

	// publish service
	publishGroup := h.Group("/publish")
	publishGroup.POST("/action", PublishAction)
	publishGroup.GET("/list")

	// favorite service
	favoriteGroup := h.Group("/favorite")
	favoriteGroup.POST("/action")
	favoriteGroup.GET("/list")

	// comment service
	commentGroup := h.Group("/comment")
	commentGroup.POST("/action")
	commentGroup.GET("/list")

	// relation service
	relationGroup := h.Group("/relation")
	relationGroup.GET("/follow/list")
	relationGroup.GET("/follower/list")
	relationGroup.GET("/friend/list")

	// message service
	messageGroup := h.Group("/message")
	messageGroup.POST("/action")
	messageGroup.GET("/chat")

	url := swagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, url))

	h.Spin()
}
