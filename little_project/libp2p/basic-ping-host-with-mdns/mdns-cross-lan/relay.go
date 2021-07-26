package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"io/ioutil"
	"khalid.fondation/libp2pdemo/basic-ping-host-with-mdns/msdn"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go

func main() {
	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	certBytes, err := ioutil.ReadFile("./server.key")
	if err != nil {
		log.Println("unable to read client.pem")
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	//priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("unable to unmarshal, err: ", err)
		return
	}
	log.Println("priv: ", priv)

	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/22001"),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.Ping(false),
		libp2p.Identity(priv),
		libp2p.EnableNATService())

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

	peerChan := msdn.InitMDNS(ctx, node, "meetup-13")

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
