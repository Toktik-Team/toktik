package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"toktik/config"
	auth "toktik/kitex_gen/douyin/auth/authservice"
)

func main() {
	var err error

	r, err := consul.NewConsulRegister(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	srv := auth.NewServer(new(AuthServiceImpl), server.WithRegistry(r), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: config.AuthServiceName,
	}))
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
