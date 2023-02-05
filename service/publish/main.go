package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"toktik/config"
	publish "toktik/kitex_gen/douyin/publish/publishservice"
	"toktik/service/publish/storage"
)

func main() {
	storage.Init()
	var err error

	addr, err := net.ResolveTCPAddr("tcp", config.PublishServiceAddr)
	if err != nil {
		panic(err)
	}
	svr := publish.NewServer(new(PublishServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
