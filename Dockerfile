# 编译镜像
FROM docker.io/golang:1.19.5-bullseye AS build
ENV TZ=Asia/Shanghai

# 配置 git
RUN apt update
RUN apt install git

# 获取文件
RUN mkdir -p /source
WORKDIR /source
COPY . .

# 编译
RUN bash build-all.sh

# 运行环境
FROM gcr.io/distroless/base-debian11
ENV TZ=Asia/Shanghai

# RUN mkdir -p /data/apps/nico-minidouyin-web
WORKDIR /data/apps/nico-minidouyin-web

# 收集数据
COPY --from=build /source/start.sh .
COPY --from=build /source/output/ .
ENTRYPOINT ["./start.sh"]
