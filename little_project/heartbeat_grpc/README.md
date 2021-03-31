### 项目描述
基于gRPC的双向流式RPC实现双向通讯应用，主要是server定期反馈在线的agent，agent定期向server上报心跳，并使用长连接来做

### 实现
主要参考：https://grpc.io/docs/languages/go/basics/
下载 protoc: https://github.com/protocolbuffers/protobuf/releases/tag/v3.15.6
环境准备如下：
```shell
export GO111MODULE=on
go get google.golang.org/protobuf/cmd/protoc-gen-go \ google.golang.org/grpc/cmd/protoc-gen-go-grpc
```
准备好proto文件之后，执行如下命令
```shell
protoc --go_out=. --go_out=paths=source_relative --go-grpc_out=. --go_out=paths=source_relative \
	proto/stream.proto
```
