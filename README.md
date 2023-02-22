# Toktik

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | A simple short video software microservice backend with basic media stream functions and social functions support, made by [Toktik-Team](https://github.com/Toktik-Team) in *The 5th Bytedance Youth Training Camp*. |
|----------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

**English** | [简体中文](README_zh-CN.md)

## Demo（With Backend API Only）

https://toktik.xctra.cn/

## Project Structure

- [constant](constant)
    - [biz](constant/biz) - Business logic related constants
    - [config](constant/config)
        - [env.go](constant/config/env.go) - Environment variables
        - [service.go](constant/config/service.go) - Service name and port
- [idl](idl)
    - [auth.proto](idl/auth.proto) - RPC definition of auth service
    - [comment.proto](idl/comment.proto) - RPC definition of comment service
    - [favorite.proto](idl/favorite.proto) - RPC definition of favorite service
    - [feed.proto](idl/feed.proto) - RPC definition of feed service
    - [publish.proto](idl/publish.proto) - RPC definition of publish service
    - [relation.proto](idl/relation.proto) - RPC definition of relation service
    - [user.proto](idl/user.proto) - RPC definition of user service
    - [wechat.proto](idl/wechat.proto) - RPC definition of chat service
- [kitex_gen](kitex_gen) - Generated code by kitex
- [logging](logging) - Logging middleware
- [manifests-dev](manifests-dev) - Kubernetes manifests for development
- [repo](repo) - Database schema and generated code by gorm gen
- [service](service)
    - [auth](service/auth) - Auth service impl
    - [comment](service/comment) - Comment service impl
    - [favorite](service/favorite) - Favorite service impl
    - [feed](service/feed) - Feed service impl
    - [publish](service/publish) - Publish service impl
    - [relation](service/relation) - Relation service impl
    - [user](service/user) - User service impl
    - [web](service/web) - Web gateway with middleware support
    - [wechat](service/wechat) - Chat service impl
- [storage](storage) - Storage middleware, s3 and local storage supported
- [test](test)
    - [e2e](test/e2e) - End-to-end test
    - [mock](test/mock) - Mock data for unit test

## Compile

1. Run the following command to install the required tools:

```bash
apt-get update && \
    apt-get install -yq git ffmpeg libavcodec-dev libavutil-dev libavformat-dev libswscale-dev && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*
    
brew install protobuf
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

2. Run `./build-all.sh` in your Linux environment to compile all services.

## Configurations

- Go `./constant/config/env.go` to see the environment variables.
- The following software/service is required to run the project:
    - PostgreSQL
    - Redis
    - HashiCorp Consul
    - FFMPEG
- The following software/service is optional:
    - A web server, such as Nginx.
    - Amazon S3, or local storage.
    - OpenTelemetry
    - Jaeger
    - Victoriametrics
    - Grafana

## Run

- Run `start.sh --service <service_name>` to start a service.
- See `./service` to get the service name.

## Test

### For Unit Test

- Run `./unit-test.sh`

### For End-to-End Test

Compile and run the following files:

- `./test/e2e/base_api_test.go` for basic API test
- `./test/e2e/interact_api_test.go` for interaction API test
- `./test/e2e/social_api_test.go` for social API test

## How to Contribute

1. Please following the [HowToAddServices](docs/HowToAddServices.md) to create your own service.
2. Create a new branch and make your changes.
3. Create a pull request to the `main` branch.
4. Wait for review and merge.

## Contributors

- [Nico](https://github.com/nicognaW)
- [Eric_Lian](https://github.com/ExerciseBook)
- [Dark Litss](https://github.com/lss233)
- [HikariLan](https://github.com/shaokeyibb)
- [YunShu](https://github.com/Selflocking)

## License

Toktik is licensed under the [MIT License](LICENSE).

