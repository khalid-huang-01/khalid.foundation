### CA系统
抽取kubeedge里面的ca系统实现的证书申请系统，客户端申请的证书可以用与server与client认证棘突可以看httpserver/server.go里面的signCerts的Usage

### 使用
```shell
go run server.go
#运行client
cd acl
go test manager_test.go
```
