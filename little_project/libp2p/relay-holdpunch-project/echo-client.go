package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
)

func main()  {
	ctx := context.Background()

	node, err := libp2p.New(ctx, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	addrStr := "/ip4/127.0.0.1/tcp/10001/p2p/QmRxm1pzaBhUpvBBNCUF13GpfndbmFzPNqoUjh4SHcxSEH"
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
	}}
