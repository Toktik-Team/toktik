# Toktik

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | 一个简单的短视频软件微服务后端，具有基本的媒体流功能和社交功能支持，由 [Toktik-Team](https://github.com/Toktik-Team) 在*第五届字节跳动青少年训练营*中开发。 |
|----------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------|

[English](README.md) | **简体中文**

## 示例（仅后端 API）

https://toktik.xctra.cn/

## 项目结构
- [config](config)
    - [config.go](config/config.go) - 环境变量
- [idl](idl)
    - [auth.proto](idl/auth.proto) - 鉴权服务 RPC 定义
    - [feed.proto](idl/feed.proto) - 视频流服务 RPC 定义
    - [user.proto](idl/user.proto) - 用户服务 RPC 定义
    - [publish.proto](idl/publish.proto) - 视频发布服务 RPC 定义
- [kitex_gen](kitex_gen) - 由 kitex 自动生成的代码
- [manifests-dev](manifests-dev) - Kubernetes 清单文件
- [repo](repo) - 数据库概要和由 gorm gen 自动生成的代码
- [service](service)
    - [auth](service/auth) - 鉴权服务实现
    - [feed](service/feed) - 视频流服务实现
    - [user](service/user) - 用户服务实现
    - [publish](service/publish) - 视频发布服务实现
    - [web](service/web) - 带中间件的 web 网关服务

## 编译

- 在你的 Linux 环境中运行 `./build-all.sh` 来编译所有的服务

## 安装

- TODO

## 运行

- TODO

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

