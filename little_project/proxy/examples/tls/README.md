## 步骤
1. 服务器端的证书生成
```shell
# 生成服务器端的私钥
openssl genrsa -out server-1.key 2048
# 生成服务器端证书
openssl req -new -x509 -key server-1.key -out server.pem -days 3650
```

2. 客户端的证书生成
```shell
# 生成客户端的私钥
openssl genrsa -out client.key 2048
# 生成客户端的证书
openssl req -new -x509 -key client.key -out client.pem -days 3650
```

3. 


## 参考地址
https://colobu.com/2016/06/07/simple-golang-tls-examples/