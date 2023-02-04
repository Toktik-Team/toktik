package main

import (
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
	"toktik/config"
	auth "toktik/kitex_gen/douyin/auth/authservice"
)

func main() {
	var err error

	addr, err := net.ResolveTCPAddr("tcp", config.AuthServiceAddr)
	if err != nil {
		panic(err)
	}
	svr := auth.NewServer(new(AuthServiceImpl), server.WithServiceAddr(addr))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
