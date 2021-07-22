package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
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
	peerChan := msdn.InitMDNS(ctx, node, "meetup")


	fmt.Println("start listener")
	for peer := range peerChan {
		fmt.Println("Found peer:", peer, ", connecting")
		if err := node.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			ch := pingService.Ping(ctx, peer.ID)
			go func(ch <-chan ping.Result) {
				for i := 0; i < 5; i++ {
					res := <-ch
					fmt.Println("pinged", peer.ID, "in", res.RTT)
					time.Sleep(10 * time.Second)
				}
			}(ch)
		}
	}

	fmt.Println("exit")
	if err := node.Close(); err != nil {
		panic(err)
	}
}
