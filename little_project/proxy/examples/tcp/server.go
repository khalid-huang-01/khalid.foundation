package main

import (
	"log"
	"net"
	"strings"
)

// 监听本地30002端口，实现echo服务
func main() {
	port := "30002"
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("err: ", err)
	}
	log.Println("Listening on TCP port: ", port)
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("err: ", err)
		}
		go handleTCPRequest(conn)
	}
}

func handleTCPRequest(conn net.Conn) {
	defer conn.Close()
	// 写一段语
	message := "echo service"
	conn.Write([]byte(message))
	// 不断回写
	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			log.Println("err: ", err)
			return
		}
		data := buf[:size]
		if strings.TrimSpace(string(data)) == "STOP" {
			log.Println("Exiting TCP server")
			return
		}
		log.Printf("Received Data: %s", data)
		conn.Write(data)
	}
}
