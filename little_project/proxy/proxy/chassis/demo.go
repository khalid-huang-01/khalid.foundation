package main

import (
	"log"
	"net"

	"khalid.foundation/proxy/proxy/chassis/protocol"
)

func main() {
	port := "30001"
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("err:", err)
		return
	}
	log.Println("Listening on TCP port: ", port)
	// 开始代理请求
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	p := &protocol.TCP{Conn: conn}
	p.Process()
}
