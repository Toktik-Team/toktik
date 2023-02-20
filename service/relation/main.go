package main

import (
	"log"
	"net"
	"toktik/constant/config"
	relation "toktik/kitex_gen/douyin/relation/relationservice"

	"github.com/cloudwego/kitex/server"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	consul "github.com/kitex-contrib/registry-consul"
)

func main() {
	r, err := consul.NewConsulRegister(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", config.RelationServiceAddr)
	if err != nil {
		panic(err)
	}

	srv := relation.NewServer(
		new(RelationServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.RelationServiceName,
		}),
	)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
