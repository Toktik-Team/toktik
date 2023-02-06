package main

import (
	"log"
	"net"
	"toktik/config"
	publish "toktik/kitex_gen/douyin/publish/publishservice"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {
	r, err := consul.NewConsulRegister(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", config.PublishServiceAddr)
	if err != nil {
		panic(err)
	}

	srv := publish.NewServer(
		new(PublishServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.PublishServiceName,
		}),
	)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
