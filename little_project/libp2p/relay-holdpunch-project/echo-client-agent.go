package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
)

func main()  {
	ctx := context.Background()

	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.EnableRelay(circuit.OptActive),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10003"))
	if err != nil {
		panic(err)
	}

	// 拼凑出一个基于relay访问echo-client-agent的地址，{relay-agent-addr}/p2p-circuit/p2p/{echo-server-agent-id}
	addrStr := "/ip4/192.168.0.10/tcp/10001/p2p/Qma57EzsaP5FyNpqduXpCgR9xUkWoXkkKoZurTH9rqPBgC/p2p-circuit/p2p/QmQxN4yj146YyXgEfYXD57HTPzmHugnXVuUAHE4HJ9FKzU"
	//addrStr := "/ip4/192.168.0.10/tcp/10002/p2p/QmQxN4yj146YyXgEfYXD57HTPzmHugnXVuUAHE4HJ9FKzU"
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

	conns := node.Network().ConnsToPeer(peer.ID)
	for _, conn := range conns {
		fmt.Printf("conn: %v\n", conn)
	}

	s, err := node.NewStream(ctx, peer.ID, "/echo/1.0.0")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sender syaning hello")
	_, err = s.Write([]byte("Hello, world\n"))

	if err != nil {
		log.Println(err)
		return
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("read reply: %q\n", out)

	if err := node.Close(); err != nil {
		panic(err)
	}
}
