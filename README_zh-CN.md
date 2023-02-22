# Toktik

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | 一个简单的短视频软件微服务后端，具有基本的媒体流功能和社交功能支持，由 [Toktik-Team](https://github.com/Toktik-Team) 在*第五届字节跳动青少年训练营*中开发。 |
|----------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|

[English](README.md) | **简体中文**

## 示例（仅后端 API）

https://toktik.xctra.cn/

## 项目结构

- [constant](constant)
  - [biz](constant/biz) - 业务逻辑相关常量
  - [config](constant/config)
    - [env.go](constant/config/env.go) - 环境变量
    - [service.go](constant/config/service.go) - 服务名称和端口
- [idl](idl)
  - [auth.proto](idl/auth.proto) - 鉴权服务 RPC 定义
  - [comment.proto](idl/comment.proto) - 评论服务 RPC 定义
  - [favorite.proto](idl/favorite.proto) - 点赞服务 RPC 定义
  - [feed.proto](idl/feed.proto) - 视频流服务 RPC 定义
  - [publish.proto](idl/publish.proto) - 视频发布服务 RPC 定义
  - [relation.proto](idl/relation.proto) - 关注服务 RPC 定义
  - [user.proto](idl/user.proto) - 用户服务 RPC 定义
  - [wechat.proto](idl/wechat.proto) - 聊天服务 RPC 定义
- [kitex_gen](kitex_gen) - 由 kitex 自动生成的代码
- [logging](logging) - 日志中间件配置
- [manifests-dev](manifests-dev) - Kubernetes 清单文件
- [repo](repo) - 数据库概要和由 gorm gen 自动生成的代码
- [service](service)
  - [auth](service/auth) - 鉴权服务实现
  - [comment](service/comment) - 评论服务实现
  - [favorite](service/favorite) - 点赞服务实现
  - [feed](service/feed) - 视频流服务实现
  - [publish](service/publish) - 视频发布服务实现
  - [relation](service/relation) - 关注服务实现
  - [user](service/user) - 用户服务实现
  - [web](service/web) - 带中间件的 web 服务实现
  - [wechat](service/wechat) - 聊天服务实现
- [storage](storage) - 对象存储中间件，支持 Amazon S3 和本地存储
- [test](test)
  - [e2e](test/e2e) - 端到端测试
  - [mock](test/mock) - 用于单元测试的 mock 数据

## 编译

1. 运行以下命令以安装所需软件包:

```bash
apt-get update && \
    apt-get install -yq git ffmpeg libavcodec-dev libavutil-dev libavformat-dev libswscale-dev && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*
    
brew install protobuf
go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
```

2. 运行 `./build-all.sh` 以编译所有服务.

## 配置

- 前往 `./constant/config/env.go` 以查看所需环境变量.
- 以下软件/服务是运行 toktik 所必需的:
  - PostgreSQL
  - Redis
  - HashiCorp Consul
  - FFMPEG
- 以下软件/服务则是可选的:
  - 一个 Web 服务器，例如 Nginx.
  - Amazon S3, 否则使用本地存储.
  - OpenTelemetry
  - Jaeger
  - Victoriametrics
  - Grafana

## 运行

- 运行 `start.sh --service <service_name>` 来启动指定服务.
- 查看 `./service` 文件夹以获得服务列表.

## 测试

### 对于单元测试

- 运行 `./unit-test.sh`

### 对于端到端测试

编译并运行以下文件来进行端到端测试:

- `./test/e2e/base_api_test.go` 对于基本 API 测试
- `./test/e2e/interact_api_test.go` 对于交互 API 测试
- `./test/e2e/social_api_test.go` 对于社交 API 测试

## 如何贡献

1. 请遵循 [HowToAddServices](docs/HowToAddServices.md) 文件说明以创建新的服务。
2. 创建一个新的分支并做出更改。
3. 提交一个 Pull Request 到 `main` 分支。
4. 等待 review 和合并。

## 贡献者

- [Nico](https://github.com/nicognaW)
- [Eric_Lian](https://github.com/ExerciseBook)
- [Dark Litss](https://github.com/lss233)
- [HikariLan](https://github.com/shaokeyibb)
- [YunShu](https://github.com/Selflocking)

## 协议

Toktik is licensed under the [MIT License](LICENSE).

