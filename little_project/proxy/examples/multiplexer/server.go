package main

import (
	"fmt"
	"log"
	"net"

	"github.com/libp2p/go-yamux/v2"
)

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
		session, err := yamux.Server(conn, nil)
		if err != nil {
			panic(err)
		}
		stream, err := session.Accept()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 4)
		stream.Read(buf)
		fmt.Println(string(buf))
	}
}