https://blog.laisky.com/p/go-mutual-tls-tcp/
https://www.jianshu.com/p/14fbf585b0d3

### 生成证书等信息
```shell
# 1. 生成 CA
openssl genrsa  -out ca.key 4096
openssl req -new -x509 -days 365 -key ca.key -out ca.crt

# 2. 生成服务端秘钥
openssl genrsa -out server.key 1024

# 3. 生成服务端证书的 CSR
# CN(common name) 必须填写为: localhost
openssl req -new -key server.key -out server.csr

# 4. 通过 CSR 向 CA 签发服务端证书
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# 5. 生成客户端秘钥
openssl genrsa -out client.key 1024

# 6. 生成客户端 CSR，一路回车即可
openssl req -new -key client.key -out client.csr

# 7. 通过 CSR 向 CA 签发客户端证书
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out client.crt
```

#### 单向认证
```shell
go run ./server/server.go
go run ./client/client.go
```

### 双向认证
```shell
go run ./server/mutal_tls_server.go
go run ./server/mutal_tls_client.go
```
