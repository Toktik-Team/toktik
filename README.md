# Toktik

| ![logo](https://avatars.githubusercontent.com/u/124244470?s=200&v=4) | A simple short video software microservice backend with basic media stream functions and social functions support, made by [Toktik-Team](https://github.com/Toktik-Team) in *The 5th Bytedance Youth Training Camp*. |
|----------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|

**English** | [简体中文](README_zh-CN.md)

## Demo（With Backend API Only）

https://toktik.xctra.cn/

## Project Structure

- [config](config)
    - [config.go](config/config.go) - Environment variables
- [idl](idl)
    - [auth.proto](idl/auth.proto) - RPC definition of auth service
    - [feed.proto](idl/feed.proto) - RPC definition of feed service
    - [user.proto](idl/user.proto) - RPC definition of user service
    - [publish.proto](idl/publish.proto) - RPC definition of publish service
- [kitex_gen](kitex_gen) - Generated code by kitex
- [manifests-dev](manifests-dev) - Kubernetes manifests for development
- [repo](repo) - Database schema and generated code by gorm gen
- [service](service)
    - [auth](service/auth) - Auth service impl
    - [feed](service/feed) - Feed service impl
    - [user](service/user) - User service impl
    - [publish](service/publish) - Publish service impl
    - [web](service/web) - Web gateway with middleware support

## Compile

- Run `./build-all.sh` in your Linux environment to compile all services.

## Install

- TODO

## Run

- TODO

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

