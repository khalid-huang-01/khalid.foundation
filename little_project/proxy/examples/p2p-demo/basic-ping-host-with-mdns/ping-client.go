package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"khalid.foundation/proxy/msdn/msdn"
	"time"
)

func main()  {
	ctx := context.Background()

	node, err := libp2p.New(ctx, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}

	pingService := &ping.PingService{Host: node}
	peerChan := msdn.InitMDNS(ctx, node, "meetup-3")


	fmt.Println("start listener")
	for p := range peerChan {
		fmt.Println("Found p:", p, ", connecting")
		if err := node.Connect(ctx, p); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}
		fmt.Println("Connection success")
		ch := pingService.Ping(ctx, p.ID)
		go func(ch <-chan ping.Result, addrInfo peer.AddrInfo) {
			for i := 0; i < 5; i++ {
				res := <-ch
				fmt.Println("pinged", addrInfo.Addrs, "in", res.RTT)
				time.Sleep(10 * time.Second)
			}
		}(ch, p)
	}

	fmt.Println("exit")
	if err := node.Close(); err != nil {
		panic(err)
	}
}
