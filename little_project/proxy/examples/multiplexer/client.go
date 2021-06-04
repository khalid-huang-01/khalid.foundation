package main

import (
	"context"
	"net"

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
	 stream, err := session.Open(context.Background())
	 if err != nil {
	 	panic(err)
	 }
	 stream.Write([]byte("ping"))
	select {

	 }
}