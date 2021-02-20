# 容器通讯场景

## 主机容器间通信
### 实践
```shell
docker run -it busybox # 开终端1，创建容器busybox-1
docker run -it busybox # 开终端2，创建容器busybox-2
ifconfig # 在busybox-1和busybox-2中查看网址，假设busybox-1为172.0.0.2, busybox-2为172.0.0.3
ping 172.0.0.3 # 在busybox-1容器中执行ping命令，发现可以互通
```
### 原理
上面就是经典容器组网模型 veth pair + bridge的模式  
在主机上安装完Docker后，Docker会在宿主机上默认创建一个名叫docker0的网桥(可以通过在主机执行`brctl show` 查看)，当我们创建一个容器的时候，如果没有指定组网模型，默认的就是使用bridge模式，在这种模式下创建容器，同时会有一对veth pair的虚拟设备被创建出来，其中一端作为容器的虚拟网络设备（也就是我们通过ifconfig看到的`eth0`），另一端会连接在docker0网桥上。那么这有什么用呢？  
Veth是虚拟以太网卡（Virtual Ethernet）的缩写，Veth Pair设备的特点就是，它被创建出来后，会以成对的虚拟网卡出现。veth pair一端发送的数据会在另一端接收，哪怕这两端被分配到不同的Network Namespace里
那根据上面 对veth pair的解读，我们就可以知道，我们的网络拓扑图如下（《极客时间里面的图》）
我们在busybox-1调用ping 172.0.0.3的时候，数据先到虚拟网卡eth0，然后直接到docker0网桥，再由docker0网桥发送到busybox-2的虚拟网卡eth0
所以说，借助网桥docker0，就可以实现一个主机上容器A和B的直接通信

### 补充
#### 容器与host veth pair的关系 
在创建容器之后调用ifconfig，我们只可以看到一个eth0，但是我们怎么知道主机上的vethxxx和哪个container eth0是成对的呢？或者说我们是否可以验证，container eth0真的有一个对端在主机上呢？
```shell
# 在目标容器中，查看eth0的对端index，在/sys/class/net/eth0下面有一个ifindex表示这个端对应的Index,而iflink表示 的是对端的index
cat /sys/class/net/eth0/iflink
# 在主机查看，确认vethxxx的index与上面的iflink是一致的（建议只创建一个容器，方便验证）
cat /sys/class/net/vethxxxx/ifindex
```
#### Docker四大网络模式[参考《Kubernetes网络权威指南》]
1. network namespace + bridge模式：Docker在安装时会创建一个名为docker0的Linux网桥，bridge模式是Docker默认的网络模式，在这种模式下，Docker会为每一个容器分配 network namespace、设置IP等，并将Docker容器连接到docker0网桥上，严谨的说，创建的容器的veth pair的一端桥接到docker0上；bridge模式为Docker容器创建独立的网络栈，保证容器内的进程使用独立的网络环境，使容器和容器、容器和宿主机之间能实现网络隔离，并可以通过网桥进行交互
2. host模式: 连接到host网络的容器共享Docker host的网络栈，容器的网络配置与host完全一样。host模式下容器将不会获得独立的network namespace，而不是和宿主机共用 一个network namespace。容器将不会虚拟出自己的网卡，配置自己的IP，而是使用宿主机的IP和端口
3. container模式：共享IP和端口，创建容器时使用--network=container:NAME_OR_ID模式，在创建新的容器时指定容器的网络和一个已经存在的容器共享一个netowrk namespace，但是并不为Docker容器进行任何网络配置，这个容器没有网卡，没有IP等，需要自行配置（可以通过lo网卡设备通信）=> kubernetes的Pod网络采用的就是Docker的container模式网络
4. none模式：none模式下的容器只有lo回环网络，没有其他网卡。主要作用是可以自行配置 

## 容器与主机外网
### 实践
我们通过network namespace + bridge的方式，可以直接在容器里面去ping通容器的IP；在与主机外网做交互时，分如下两种
#### 主机外网访问容器服务
如果容器需要暴露服务给外网访问的话，一般是使用端口映射的方式，其原理是在本地的iptable的nat表中添加相应的规则，将访问IP地址：hostport的网包进行一次DNAT，转成容器IP：containerport，然后就会根据route到docker0网桥，再到对应的容器网络中
```shell
# 部署nginx并将容器端口1234映射到主机端口80
docker run -p 1234:80 -d nginx 
docker ps 
# 通过 iptable验证; 在下面的命令可以看到DNAT发生在DOCKER这条iptables链，它有两处引用，分别是PREROUTING链和OUTPUT链，意味着从外面发到本机和本地进程访问本机（由iptables匹配规则 ADDRTYPE match dst-type LOCAL指定）的1234端口的包目的地址都会被修改为172.17.0.2:80
iptables -t nat -nL
```
#### 容器访问主机外网
容器内访问外网需要ip_forward和SNAT/MASQUERADE，默认情况下都是开通的，
```shell
# 主机ip_forward打开 
net.ipv4.ip_forward=1
# Docker会自动 在iptables的POSTROUTING链上创建如下规则：
Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination         
MASQUERADE  all  --  172.17.0.0/16        0.0.0.0/0  
# 即从容器网段出来访问外网的包，都要做MASQUERADE，即出去的包都用主机的IP地址替换源地址
```

### 原理

## 跨主机容器间通信
### 原理
我们在前面看到的都是单机的容器通信，docker在一开始设计的时候就没有考虑跨主机容器通信问题，目前业界主流的跨主机容器通信解决方案大致分为“隧道方案”和“路由方案”  
#### 隧道方案
隧道网络也称为overlay网络，或者覆盖网络。这种网络的优点在于适用于几乎所有网络基础架构，它唯一的要求是节点之间是三层互通的，也就是IP互通； 但是它的问题是随着节点规模增长，复杂度也会增长，而且由于使用了封解包技术，出了网络问题也很难解决。
典型的实现有flannel支持的UDP封包及其Linux内核的VXLAN协议。具体会以k8s网络支持部分进行解读

这里不做特别的实践演示，具体的实践将在k8s的场景下进行

## k8s单pod内容器间通信
### 实践
多容器pod，我们在一个pod里面部署一个nginx容器和一个busybox容器，尝试通过busybox容器访问nginx容器，确认在单pod里面的容器通信
```yaml
apiVersion: v1
kind: pod
metadata:
    name: busybox-nginx
spec:
    containers:
    - name: busybox
      image: busybox
      args:
      - sleep
      - "12000"
    - name: nginx
      image: nginx:1.16.1
      ports:
      - containerPort: 80
```
执行如下命令：
```shell
kubectl apply -f <file-name>
# 进入容器
kubectl exec -it busybox-nginx -c busybox -- /bin/sh
# 在busybox容器中通过localhost去访问nginx的服务
wget -O - localhost:80 
```
### 原理
kubernetes是“单pod单IP”模型，也就是每个pod都有一个独立的IP，Pod内所有容器共享network namespace(同一个网络协议栈和IP)。基于“单pod单IP”模型构建的kubernetes扁平网络里面，容器是一等公民，容器之间以及容器与Node直接通信，不需要额外的NAT，减少了NAT带来的性能损耗，而且可追溯源地址。（那么pod里面的容器如何访问外网以及暴露服务呢？总体来说，就是通过Service和Ingress），那么“单pod单IP”是如何实现的呢？具体的流程如下：
1. 创建pod之后，会调用CRI创建pod内的若干容器
2. 首先创建一个pause容器，占用一个network namespace
3. 其他用户容器加入pause容器占用的network namespace，从而实现pod内容器共享同一个network namespace
所以在“单pod单IP”下pod里面的每个容器都会看到相同的网络，每个容器都通过localhost就能访问同pod下的其他容器；而且k8s的底层网络是“全连通”的，即在同一集群内运动的所有Pod都可以自由通信

## k8s 集群内pod与pod容器间通信（pod容器之间可以直达，flannld）
### 实践
部署一个带flannel的k8s集群，部署方式可以参考""  
部署一个nginx server，由service + deployment组成，然后再部署一个busybox pod，首先通过busybox 容器直接利用pod ip验证是否可以通过pod ip访问，再使用推荐的service方式进行访问; 在部署的时候，将上面两个应用部署到不同的节点上
```shell
# 创建busybox
kubectl apply -f <<EOF
apiVersion: v1
kind: Pod
metadata:
    name: busybox-sleep-resources
spec:
    nodeName: node2
    containers:
    - name: busybox
      image: busybox
      args:
      - sleep
      - "12000"
EOF
# 创建nginx-server
kubectl apply -f <<EOF
apiVersion: v1
kind: Service
metadata:
    name: nginx-server-service
spec:
    ports:
    - port: 80
      name: http
      targetPort: 80
    selector:
      app: nginx
    type: NodePort
---
apiVersion: apps/v1
kind: Deployment
metadata:
    name: nginx-server
spec:
    selector:
        matchLabels:
            app: nginx
    template:
        metadata:
            labels:
                app: nginx
        spec:
            nodeName: node1
            containers:
            - name: nginx
              image: nginx:1.16.1
              ports:
              - continerPort: 80
EOF
# 通过kubectl get pod -o wide 获取 nginx-server的pod的ip，为10.244.1.35
# 进入busybox容器中
kubectl exec -it busybox-sleep-resources -- /bin/sh
# 在容器中执行
wget -O - 10.244.1.35:80 # 正常情况下可互通，访问到页面
```
### 原理
通过前面对“主机容器间通信”的实践，我们知道，通过docker0网桥我们可以实现一个主机上容器A和B的直接通信。如果需要实现位于不同宿主机上的容器之间的通信，就需要使用到类似flannel的容器网络解决方案。flannel支持UDP模型和VXLAN模式，相比于UDP模型VXLAN模式会更高效。具体的内容可以参考《Kubernetes网络权威指南》，这里根据书本的内容对VXLAN模型做下说明
XVLAN（virtual eXtensible LAN，虚拟可扩展的局域网）是一种虚拟化隧道通信技术。它是一种overlay技术，通过三层的网络搭建虚拟的二层网络（二层网络仅通过MAX寻址即可实现通讯，一般是小局域网；但三层网络通过IP路由实现跨网段通讯，三层网络则可以组大型网络）
kubernetes当前只支持单网络，所以在三层网络上只有一个VXLAN网络，flannel会在集群的节点上创建一个名为flannel.1的VXLAN网卡（区别于UDP模式下的flannel0是一个tun设备）
数据路径如下：
[图片]()  
1. 容器A中的IP包通过容器A内的路由表被发送到cnio
2. 到达cnio中的IP包通过匹配host A中的路由表发现通往10.244.2.194的IP包应该交给flannel.1接口
3. flannel.1作为VTEP设备，收到报文后将按照VTEP的配置进行封包。并行通过 etcd得知10.244.2.194属于节点B，并得到节点B的IP，然后通过节点A中的转发表得到节点B对应的VTEP的MAC，根据flannel.1设备创建时的设置参数（VNI、local IP 、port） 进行VXLAN封包
4. 通过host A跟host B之间的网络连接，VXLAN包到达host B的eth1接口
5. 通过端口8472，VXLAN包被转发给VTEP设备flannel.1进行解包
6. 解封装后的IP包匹配host B中的路由表（10.244.2.0）,内核将IP包转发给cni0
7. cni0将IP包转发给连接在cni0上的容器B  
在VXLAN模式下，数据是由内核转发的，flannel不转发数据，仅动态配置ARP和FDB表项

## k8s 集群内访问pod容器服务（service）
### 虚拟IP实践
k8s之所以会需要用于service，主要有两层原因，一个是Pod的IP是不固定的，另一个是为一组Pod实例之间总会有负载均衡的需求。
继续“k8s 集群内pod与pod容器间通信”的实验
```shell
# 获取service的虚拟IP为10.101.78.74，端口为80
kubectl get svc nginx-server-service
# 在busybox的容器里
wget -O - 10.101.78.74:80 # 正常是可以访问到nginx服务的
```
### 原理（内容主要来自张磊大佬的文章）
这里主要陈述下service的实现原理，Service是由kube-proxy组件加上iptables共同实现的。在前面我们创建一个nginx-server-service的service，kube-proxy会通过service的informer获取到这个创建信息，并在宿主机上创建一条iptables规则，具体查看如下：
```shell
iptables-save -t nat #查到如下的一条信息
-A KUBE_SERVICES -d 10.101.78.74/32 -p tcp -m comment --comment "default/nginx-server-service:http cluster IP" -m tcp --dport 80 -j KUBE-SVC-5MSM3UZHD6Y54CKX
```
这条iptable的规则 含义是：凡是目的地址是10.101.78.74、目的端口是80的IP包，都跳转到另外一条名叫KUBE-SVC-5MSM3UZHD6Y54CKX5的iptables链处理。10.101.78.74是service的IP，而80是service的port. 所以我们就可以知道service的IP其实都是VIP，它只是一条iptables规则上的配置，所以ping它是不会有任何响应的。我们继续上面的iptables链
```shell
# 为了更好的展示，我们这里把nginx-server的replicas设置为2，重新部署下后，执行
iptables-save -t nat | grep KUBE-SVC-5MSM3UZHD6Y54CKX
# 可以得到
-A KUBE-SVC-5MSM3UZHD6Y54CKX -m comment --comment "default/nginx-server-service:http" -m statistic --mode random --probability 0.50000000 -j KUBE-SEP-E6XQDMRJWOVB7TWM
-A KUBE-SVC-5MSM3UZHD6Y54CKX -m comment --comment "default/nginx-server-service:http" -j KUBE_SEP-PD20DXQ6B5WEYRQL
```
我们可以看到上面的规则 就是一组随机模式的iptables链，转发的目的是KUBE-SEP-E6XQDMRJWOVB7TWM和KUBE_SEP-PD20DXQ6B5WEYRQL，我们看下KUBE_SEP-PD20DXQ6B5WEYRQL的内容如下
```shell
-A KUBE-SEP-PD20DXQ6B5WEYRQL -s 10.244.1.36/32 -m comment --comment "default/nginx-server-service:http" -j KUBE-MARK-MASQ
-A KUBE-SEP-PD20DXQ6B5WEYRQL -p tcp -m comment --comment "default/nginx-server-service:http" -m tcp -j DNAT --to-destination 10.244.1.36:80
```
主要是DNAT规则 ，就是在PREROUTING检查点之前，将注入IP包的目的地址和端口，改成-to-destination所指定 的新的目的地址和端口，这个目的地址和端口就是被代理的Pod的IP地址和端口  

当然这种方法方法会随着POD数量的上升导致 iptables数量上升，占用大量宿主机CPU资源，当前的解决方案是使用IPVS模式的Service。这里不做介绍了，IPVS主要是将规则的处理放到内核态，降低了维护这些规则 的代价

### 域名实践
```shell
wget -O - nginx-server-service.default # <service-name>.<namespace-name> 正常这个也是可以访问到nginx服务的
```
### 原理（内容主要）
这里我们主要知道 Service和Pod都会被分配对应的DNS A记录（从域名解析IP的记录）

## k8s 集群外访问pod容器服务（ingress）
### 实践
segmentfault.com/a/1190000019908991部署
首先我们需要安装一下ingress-controller，
```shell
# 创建ingress-controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.20.0/deploy/mandatory.yaml
# 将ingress服务暴露出来
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.20.0/deploy/provider/baremetal/service-nodeport.yaml

### 原理
后续补充
