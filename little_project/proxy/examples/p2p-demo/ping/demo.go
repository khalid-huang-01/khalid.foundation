package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"log"
	"time"
)

func main()  {
	h2, err := libp2p.New(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println(h2.Addrs())

	go func() {
		time.Sleep(2 * time.Second)
		// 禁止别人来ping自己
		h1, err := libp2p.New(context.Background(), libp2p.Ping(false))
		if err != nil {
			panic(err)
		}
		h1.Connect(context.Background(), peer.AddrInfo{
			ID: h2.ID(),
			Addrs: h2.Addrs(),
		})
		fmt.Println(h1.ID())
		// 启动一个ping客户端
		ps1 := ping.NewPingService(h1)
		pctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// 去ping h2
		ts := ps1.Ping(pctx, h2.ID())
		for i := 0; i < 5; i++ {
			select {
			case res := <-ts:
				if res.Error != nil {
					log.Fatal("ping failed: ", res.Error)
				}
				log.Println("h1 ping h2 took: ", res.RTT)
			case <-time.After(time.Second * 4):
				log.Fatal("failed to receive ping")
			}
		}
	}()
	select {}


}