### 说明
1. 这个例子主要是在同个局域网里面使用mdns技术来作节点的自动发现，我们可以实现在不显示注入节点的peerID信息的情况下，直接获取到节点的信息，去ping节点
2. 运行如下：
```shell
# 启动一个中继节点（ipfs里面的boostrap节点）
go run enable-ping-node.go

# 在不输入上面节点信息的情况下，可以下面的节点可以直接获取上面的节点信息，并启动ping
go run ping-client.go 

# 在一个节点上，换个端口启动，查看是否可以支持ping 通
go run ping-client.go 

# 在另一个子网里面启动一个节点，看看是否可以ping 通
go run ping-client.go

```