package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()

	// 因为没有对应的协议做Secrutiy
	node, err := libp2p.New(ctx, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	addrStr := "/ip4/119.8.58.38/tcp/10004/p2p/QmTEZVRJYs3fSo1CXGztYjHayTdk66iESvhsHR1x7eDZmK"
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peer, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}
	fmt.Println("success connect to ", addrStr)

	//s, err := node.NewStream(ctx, peer.ID, "/echo/1.0.0")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//log.Println("Sender syaning hello")
	//_, err = s.Write([]byte("Hello, world\n"))
	//
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//out, err := ioutil.ReadAll(s)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//log.Printf("read reply: %q\n", out)
	//
	//if err := node.Close(); err != nil {
	//	panic(err)
	//}}
	select {
	}
}
