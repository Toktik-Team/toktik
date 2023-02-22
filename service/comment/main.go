package main

import (
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
	"net"
	"toktik/constant/config"
	comment "toktik/kitex_gen/douyin/comment/commentservice"
)

func main() {
	var err error

	r, err := consul.NewConsulRegister(config.EnvConfig.CONSUL_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveTCPAddr("tcp", config.CommentServiceAddr)
	if err != nil {
		panic(err)
	}

	provider.NewOpenTelemetryProvider(
		provider.WithServiceName(config.CommentServiceName),
		provider.WithExportEndpoint(config.EnvConfig.EXPORT_ENDPOINT),
		provider.WithInsecure(),
	)

	srv := comment.NewServer(
		new(CommentServiceImpl),
		server.WithServiceAddr(addr),
		server.WithRegistry(r),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: config.CommentServiceName,
		}),
	)

	err = srv.Run()

	if err != nil {
		log.Fatal(err)
	}
}
