FROM docker.io/golang:1.19.5-bullseye

# FFmpeg 及依赖
RUN apt-get update && \
    apt-get install -yq ffmpeg && \
    apt-get clean && \
    apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false && \
    rm -rf /var/lib/apt/lists/*
