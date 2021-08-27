https://blog.laisky.com/p/go-mutual-tls-tcp/
https://www.jianshu.com/p/14fbf585b0d3
cloud.tencent.com/developer/article/1640156

### 生成证书等信息
```shell
# 1. 生成 CA
openssl genrsa  -out ca.key 4096
openssl req -new -x509 -days 365 -key ca.key -out ca.crt

# 2. 生成服务端秘钥
openssl genrsa -out server.key 2048

# 3. 生成服务端证书的 CSR
# CN(common name) 必须填写
# 一般为网站域名，cnblogs.com/iiiiher/p/8085698.html
openssl req -new -key server.key -subj "/CN=127.0.0.1" -out server.csr

# 4. 通过 CSR 向 CA 签发服务端证书, 地址要填写，不然要配置InsecureSkipVerify才可以访问, 对应与SANs
#echo subjectAltName = IP:127.0.0.1 > etfile.cnf
cat >extfile.cnf<<EOF
subjectAltName=@alt_names
[alt_names]
DNS.1 = www.my.com
DNS.2 = www.alone.com
IP.1 = 192.168.0.38
IP.2 = 127.0.0.1
EOF

openssl x509 -req -days 365 -in server.csr -CA ../ca.crt -CAkey ../ca.key -set_serial 01 -out server.crt -extfile extfile.cnf
# 5. 生成客户端秘钥
openssl genrsa -out client.key 2048

# 6. 生成客户端 CSR，一路回车即可
openssl req -new -key client.key  -out client.csr

# 7. 通过 CSR 向 CA 签发客户端证书
openssl x509 -req -days 365 -in client.csr -CA ../ca.crt -CAkey ../ca.key -set_serial 01 -out client.crt
```

#### 单向认证
```shell
go run ./server/server.go
go run ./client/client.go
```
为了验证ca证书的有效性，可以把client里面的ca.crt修改为rootCA.crt，可以看到，如果不是ca签署的话，是过不了的; 
ca.crt是根证书，而rootCA.crt是自己伪造的其他的证书

### 双向认证
```shell
go run ./server/mutal_tls_server.go
go run ./server/mutal_tls_client.go
```


