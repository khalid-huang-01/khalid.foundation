## k8s service基于域名访问的服务发现

### 准备
busybox准备如下：
```yaml
apiVersion: v1
kind: Pod
	name: busybox-sleep-resource-cloud
spec:
	nodeName: ke-cloud
	containers:
		- name: busybox
		  image: busybox:latest
		  args:
		  	- sleep
		  	- "12000"
```

### HTTP 协议的域名访问
httpserver服务如下：
```
apiVersion: v1
kind: Pod
metadata:
	name: hostname-cloud
	labels:
		app: hostname-cloud
spec:
	nodeName: ke-cloud
	containers:
		- name: hostname
		  image: k8s.gcr.io/server_hostname: latest
		  imagePullPolicy: IfNotPresent
		  ports:
		  	- containerPort: 9376
---
apiVersion: v1
kind: Service
metadata:
	name: hostname-svc
spec:
	selector:
		app: hostname-cloud
	ports:
		- name: http-1
		  port: 12348
		  protocol: TCP
		  targetPort: 9376
```
进行验证
```shell
kubectl exec -it busybox-sleep-resources-cloud -- /bin/sh
wget -qO- hostname-svc.default:12348
# 或者使用如下命令
kubectl exec -it busybox-sleep-resources-cloud -- wget -qO- hostname-svc.default:12348
```

### TCP协议的域名访问实验
tcp服务具体使用：https://github.com/cjimti/go-echo
服务如下：
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tcp-echo-deployment
  labels:
    app: tcp-echo
    system: example
spec:
  replicas: 2
  selector:
    matchLabels:
      app: tcp-echo
  template:
    metadata:
      labels:
        app: tcp-echo
        system: example
    spec:
      containers:
        - name: tcp-echo-container
          image: cjimti/go-echo:latest
          imagePullPolicy: Always
          env:
            - name: TCP_PORT
              value: "2701"
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: POD_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.podIP
            - name: SERVICE_ACCOUNT
              valueFrom:
                fieldRef:
                  fieldPath: spec.serviceAccountName
          ports:
            - name: tcp-echo-port
              containerPort: 2701
---
apiVersion: v1
kind: Service
metadata:
  name: "tcp-echo-service"
  labels:
    app: tcp-echo
spec:
  selector:
    app: "tcp-echo"
  ports:
    - protocol: "TCP"
      port: 2701
      targetPort: 2701
```
执行测试：
```
kubectl exec -it busybox-sleep-resource-cloud -- telnet tcp-echo-service.default 2701
```

### websocket协议的域名访问实验
websocket与tcp服务类似，在websocket客户端用例 ws-svc.default:12345访问就可以进行连接了
具体服务可以使用：github.com/gorilla/websocket/tree/master/examples/echo
使用go build编译出二进制后，dockerfile如下：
```
FROM ubuntu:latest
EXPOSE 8080
COPY server /home/service/server
COPY client /home/service/client
WORKDIR /home/service
ENTRYPOINT ["./server", "--addr", "0.0.0.0:8080"]
```
如果想使用alpine:3.11的话，要用alpine来打镜像，否则会无法运行，具体可以参考：github.com/kubeedge/kubeedge/blog/master/build/edge/Dockerfile
yaml如下：
```
apiVersion: v1
kind: Pod
metadata:
	name: ws-cloud
	labels:
		app: ws-cloud
spec:
	nodeName: ke-cloud
	containers:
		- name : ws
   		  image: websocket:0.1
   		  ports:
   		  	- containerPort: 8080
---
apiVersion: v1
kind: Service
metadata: 
	name: ws-svc
spec:
	selector:
		app: ws-cloud
	ports:
		- name: http-1
		  port: 12348
		  protocol: TCP
		  targetPort: 8080
```
进行测试
```
kubectl exec -it ws-cloud -- ./client --addr ws-svc.default:12348
```