package main

import (
	"khalid.foundation/proxy/proxy/tcp"
	"log"
	"net"
)

var (
	laddr = "0.0.0.0:30001"
	raddr = "localhost:30002"
)

// TCP server proxy tcp 方式
// 		启动tcp server 监听在30002端口，也就是examples/tcp里面的
//		启动了服务之后，就直接 telnet 192.168.0.10 30001

// TCP serer proxy http 方式 透传
//  	启动http server监听在30002端口，也就是examples/http里面的
// 		启动本服务后， curl localhost:30001/hello
func main() {
	// start local server proxy server
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Println("Failed to open local port to listen: ", err)
		return
	}
	log.Println("listen tcp server on addr: ", laddr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept connetion: ", err)
			continue
		}
		p := tcp.New(conn,laddr, raddr)
		go p.Start()
	}
}