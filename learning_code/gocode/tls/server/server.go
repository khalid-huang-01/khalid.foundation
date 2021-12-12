package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
)

func main() {
	cert, err := tls.LoadX509KeyPair("./server/tmp/server.crt", "./server/tmp/server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}, MaxVersion: tls.VersionTLS12}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println("err: ", err)
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		// 要知道，tls的握手验证是在这里发生的，也就是在自己的携程里面做校验的，在这里会去读取证书信息，然后
		// 做握手，具体是izai
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
