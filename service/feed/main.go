package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"toktik/config"
	feed "toktik/kitex_gen/douyin/feed/feedservice"
)

func main() {
	var err error

	addr, err := net.ResolveTCPAddr("tcp", config.FeedServiceAddr)
	if err != nil {
		panic(err)
	}

	svr := feed.NewServer(new(FeedServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
