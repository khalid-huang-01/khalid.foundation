package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"khalid.foundation/proxy/msdn/msdn"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go

func main() {
	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	node, err := libp2p.New(ctx, libp2p.Ping(false))
	if err != nil {
		panic(err)
	}
	// 这里其实就等于libp2p.New(libp2p.Ping(true))
	// 也等于pingService := ping.NewPingService(node)
	// 这里要知道pngService其实是一个server+client，可以直接用ps来ping别人的，是无状态的
	pingService := &ping.PingService{Host: node}
	node.SetStreamHandler(ping.ID, pingService.PingHandler)

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

	peerChan := msdn.InitMDNS(ctx, node, "meetup")

	for peer := range peerChan {
		fmt.Println("Found peer:", peer, ", connecting")
		if err := node.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
		}
	}


	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
	if err := node.Close(); err != nil {
		panic(err)
	}

}
