package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"os"
)

func main()  {
	ctx := context.Background()

	node, err := libp2p.New(ctx, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	if len(os.Args) <= 1 {
		panic("Please provide the peer addr")
	}
	addr, err := multiaddr.NewMultiaddr(os.Args[1])
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

	fmt.Println("sending 5 ping message to ", addr)
	// 下面也可以直接 ping.Ping(ctx, node, peer.ID)
	pingService := &ping.PingService{Host: node}
	ch := pingService.Ping(ctx, peer.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		if res.Error != nil {
			fmt.Println("ping error: ", res.Error)
			return
		}
		fmt.Println("pinged", addr, "in", res.RTT)
	}

	if err := node.Close(); err != nil {
		panic(err)
	}
}
