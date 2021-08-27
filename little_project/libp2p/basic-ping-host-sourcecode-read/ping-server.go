package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/libp2p/go-tcp-transport"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go

func main() {
	// openssl genrsa -out rsa_private.key 2048
	certBytes, err := ioutil.ReadFile("./host-acl-with-connectiongater/server.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	node, err := libp2p.New(ctx,
		libp2p.Ping(false),
		libp2p.Identity(priv),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"))
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

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
	if err := node.Close(); err != nil {
		panic(err)
	}

}
