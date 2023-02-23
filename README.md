# Toktik

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | Short video microservice application built with `Kitex` and `Hertz` , made by [Toktik-Team](https://github.com/Toktik-Team) in _The 5th Bytedance Youth Training Camp_. |
| -------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |

**English** | [简体中文](README_zh-CN.md)

## Demo（With Backend API Only）

https://toktik.xctra.cn/

## Project Structure

- [constant](constant)
  - [biz](constant/biz) - Business logic related constants
  - [config](constant/config)
    - [env.go](constant/config/env.go) - Environment variable configs
    - [service.go](constant/config/service.go) - Service names and ports
- [idl](idl)
  - [auth.proto](idl/auth.proto) - RPC definition of auth service
  - [comment.proto](idl/comment.proto) - RPC definition of comment service
  - [favorite.proto](idl/favorite.proto) - RPC definition of favorite service
  - [feed.proto](idl/feed.proto) - RPC definition of feed service
  - [publish.proto](idl/publish.proto) - RPC definition of publish service
  - [relation.proto](idl/relation.proto) - RPC definition of relation service
  - [user.proto](idl/user.proto) - RPC definition of user service
  - [wechat.proto](idl/wechat.proto) - RPC definition of chat service
- [kitex_gen](kitex_gen) - Generated code by Kitex
- [logging](logging) - Logging middleware
- [manifests-dev](manifests-dev) - Kubernetes manifests
- [repo](repo) - Database schemas and generated code by Gorm Gen
- [service](service)
  - [auth](service/auth) - Auth service impl
  - [comment](service/comment) - Comment service impl
  - [favorite](service/favorite) - Favorite service impl
  - [feed](service/feed) - Feed service impl
  - [publish](service/publish) - Publish service impl
  - [relation](service/relation) - Relation service impl
  - [user](service/user) - User service impl
  - [web](service/web) - Web api gateway
    - [mw](service/web/mw) - Hertz middlewares
  - [wechat](service/wechat) - Chat service impl
- [storage](storage) - Storage middleware, s3, local storage & volcengine [ImageX](https://www.volcengine.com/products/imagex) supported
- [test](test)
  - [e2e](test/e2e) - End-to-end test
  - [mock](test/mock) - Mock data for unit test

## Prerequisite

> This project does not support Windows, as Kitex does not support either.

- Linux / MacOS
- Go
- FFmpeg
- PostgreSQL
- Redis
- OpenTelemetry Collector

For observability infrastructures, it's recommended to use:

- Jaeger All in one
- Victoria Metrics
- Grafana

## Build

Run `./build-all.sh` in your Linux environment to compile all services.

## Configurations

Check out [`constant/config/env.go`](constant/config/env.go)

## Run

- Run `start.sh --service <service_name>` to start a service.
- `service_name` could be any of the sub directory of `./service`.

## Test

### Unit Test

Run `./unit-test.sh`

### End-to-End Test

Run `go test toktik/test/e2e -tags="e2e"`

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
