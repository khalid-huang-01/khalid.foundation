package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/libp2p/go-yamux/v2"
)

func main()  {
	conn, err := net.Dial("tcp", "localhost:30002")
	if err != nil {
		panic(err)
	}

	session, err := yamux.Client(conn, nil)
	if err != nil {
		panic(err)
	}
	for i:=0;i < 2; i++ {
		fmt.Println("i: ", i)
		stream, err := session.Open(context.Background())
		if err != nil {
			panic(err)
		}
		go startChat(stream, i)
	}
	select {

	 }
}

func startChat(stream net.Conn, i int) {
	hello := "ping" + strconv.Itoa(i)
	stream.Write([]byte(hello))
	buf := make([]byte, 10)
	stream.Read(buf)
	fmt.Println(string(buf))
}