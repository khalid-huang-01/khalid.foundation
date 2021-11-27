package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
)

func main()  {
	ctx := context.Background()

	//relayID := ""
	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.EnableRelay(circuit.OptActive))
	if err != nil {
		panic(err)
	}



	//addrStr := "/ip4/127.0.0.1/tcp/10001/p2p/QmRxm1pzaBhUpvBBNCUF13GpfndbmFzPNqoUjh4SHcxSEH"
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
