package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"time"
)

func main()  {
	ctx := context.Background()

	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.Security(libp2ptls.ID, libp2ptls.New))
	if err != nil {
		panic(err)
	}

	//if len(os.Args) <= 1 {
	//	panic("Please provide the peer addr")
	//}
	//addr, err := multiaddr.NewMultiaddr(os.Args[1])
	addrStr := "/ip4/192.168.0.10/tcp/9000/p2p/QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
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
		time.Sleep(1 * time.Second)
	}

	if err := node.Close(); err != nil {
		panic(err)
	}
}
