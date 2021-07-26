package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"khalid.fondation/libp2pdemo/basic-ping-host-with-mdns/msdn"
	"time"
)

func main()  {
	ctx := context.Background()
	//relayAddr := ""
	//addr, err := ma.NewMultiaddr(relayAddr)
	//if err != nil {
	//	panic(err)
	//}
	//raddrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	//if err != nil {
	//	panic(err)
	//}

	node, err := libp2p.New(ctx, libp2p.Ping(true),
		//libp2p.ListenAddrStrings("/ip4/192.168.0.32/tcp/10007"),
		//libp2p.EnableAutoRelay(),
		//libp2p.StaticRelays([]peer.AddrInfo{*raddrInfo}),
		//libp2p.EnableNATService())
		)

	if err != nil {
		panic(err)
	}

	pingService := &ping.PingService{Host: node}

	// pint he node's peerInfo in multiaddr format
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address: ", addrs[0])

	peerChan := msdn.InitMDNS(ctx, node, "meetup-13")


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
				if res.Error != nil {
					fmt.Println(res.Error)
					break
				}
				fmt.Println("pinged", addrInfo.Addrs, "success, in", res.RTT)
				time.Sleep(10 * time.Second)
			}
		}(ch, p)
	}

	fmt.Println("exit")
	if err := node.Close(); err != nil {
		panic(err)
	}
}
