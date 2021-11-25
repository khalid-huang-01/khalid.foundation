package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go
// 调用libp2p的库，实现一个一样的自动获取本地IP的函数

func main() {
	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10001"))

	if err != nil {
		panic(err)
	}
	for _, v := range node.Addrs() {
		fmt.Printf("%s : %v/p2p/%s\n", "Tunnel server addr", v, node.ID().Pretty())
	}

	node.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})



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
// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}