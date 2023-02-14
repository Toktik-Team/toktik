#!/bin/bash

# 检查 go 命令是否存在于 PATH 环境变量中
if ! command -v go &>/dev/null; then
  echo "错误：go 命令未在 PATH 中找到。请安装或将其添加到 PATH 中。"
  exit 1
fi

# 检查 protoc 命令是否存在于 PATH 环境变量中
if ! command -v protoc &>/dev/null; then
  echo "错误：protoc 命令未在 PATH 中找到。请安装或将其添加到 PATH 中。"
  # 检查操作系统是否为 macOS
  if [[ $(uname) == "Darwin" ]]; then
    # 如果是macOS，则检查brew是否在PATH中
    if ! command -v brew &>/dev/null; then
      echo "错误：brew 命令未在 PATH 中找到。请安装或将其添加到 PATH 中。"
      exit 1
    else
      echo "尝试安装 protoc......"
      brew install protobuf
    fi
  fi
fi

# 再次检查 protoc 命令是否存在于 PATH 环境变量中
if ! command -v protoc &>/dev/null; then
  echo "错误：protoc 命令未在 PATH 中找到，看起来安装失败了。请手动安装。"
  exit 1
fi

# 检查 kitex 命令是否存在于 PATH 环境变量中
if ! command -v kitex &>/dev/null; then
  echo "错误：kitex 命令未在 PATH 中找到，尝试安装......"
  go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
fi

# 再次检查 kitex 命令是否存在于 PATH 环境变量中
if ! command -v kitex &>/dev/null; then
  echo "错误：kitex 命令未在 PATH 中找到，看起来安装失败了。请手动安装。"
  exit 1
fi

mkdir -p kitex_gen
kitex -module "toktik" -I idl/ idl/"$1".proto

mkdir -p service/"$1"
cd service/"$1" && kitex -module "toktik" -service "$1" -use toktik/kitex_gen/ -I ../../idl/ ../../idl/"$1".proto

go mod tidy
