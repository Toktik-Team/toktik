# How To Add Services - 如何添加一个服务

1. 确定要添加的服务的名称（如：`user`, `auth`, `feed` 等）可参考 `service/web/main.go` 下的 hertz group 或 API URL
2. 编写 proto IDL 文件，并将文件命名为`{服务名称}.proto` ，放入 idl 目录
3. 打开终端 cd 到项目根目录，然后调用 `./add-kitex-service.sh {服务名称}`
4. 此时 kitex 生成的代码将被放入 `kitex_gen` 目录（无需改动）和 `service/{服务名称}`
5. 完成服务业务逻辑并妥善修改 `service/{服务名称}/handler.go`
6. 在 `service/web` web api 中添加对 RPC 服务的调用