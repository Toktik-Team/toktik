# 编译镜像
FROM docker.io/nicognaw/toktik-basic-build:v1 AS build
ENV TZ=Asia/Shanghai
ENV DEBIAN_FRONTEND=noninteractive

# 构建依赖
RUN apt-get update && \
    apt-get install -yq git ffmpeg libavcodec-dev libavutil-dev libavformat-dev libswscale-dev && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*

# 获取文件
RUN mkdir -p /source
WORKDIR /source
COPY . .

# 编译
RUN bash unit-test.sh && bash build-all.sh

# 运行环境
FROM docker.io/nicognaw/toktik-basic:v1
ENV TZ=Asia/Shanghai
ENV DEBIAN_FRONTEND=noninteractive

# FFmpeg 及依赖
RUN apt-get update && \
    apt-get install -yq ca-certificates && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*

# RUN mkdir -p /data/apps/toktik-service-bundle
WORKDIR /data/apps/toktik-service-bundle

# 收集数据
COPY --from=build /source/output/ .
