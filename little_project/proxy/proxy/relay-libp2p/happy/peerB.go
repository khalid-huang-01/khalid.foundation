package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	relayID := "QmQ6t5SKxiT1JkiXG1YnuDyZphgZmZYZGZ4itmx2BZ2rYD"
	// libp2p.ListenAddrs的作用是什么 => 启动服务，这样别人才能通过stream连接自己，这个是默认启动的，不用配置
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay())
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	relayAddr := "/ip4/127.0.0.1/tcp/10001/p2p/" + relayID
	relayAddrInfo, err := utils.Addr2info(relayAddr)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
		log.Println("Failed to connect relay, error: ", err)
		return
	}
	host.SetStreamHandler("/cats", func(s network.Stream) {
		log.Println("Meow! It worked!")
		s.Close()
	})

	utils.PrintHostAddr("Server peerB : ", host)
	log.Println("ID: ", host.ID())
	select {}
}
