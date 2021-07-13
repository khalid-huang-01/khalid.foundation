package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	relayID := "QmQuUjm1s3VfqASbRjH65MBctVQiASCLqkaMBcMj72bvtz"

	relayAddr := "/ip4/192.168.0.38/tcp/10001/p2p/" + relayID
	relayAddrInfo, err := utils.Addr2info(relayAddr)
	if err != nil {
		log.Println("err: ", err)
		return
	}
	//relayID2 := "QmU9dFLPxN9ubtb13CfgAjHRC4xMzbJQWcQWaSc8Db455i"
	//relayAddr2 := "/ip4/192.168.0.38/tcp/10001/p2p/" + relayID2
	//relayAddrInfo2, err := utils.Addr2info(relayAddr2)
	//if err != nil {
	//	log.Println("err: ", err)
	//	return
	//}


	// 也可以这样直接进行配置
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(),
		libp2p.StaticRelays([]peer.AddrInfo{*relayAddrInfo}))
		//libp2p.StaticRelays([]peer.AddrInfo{*relayAddrInfo, *relayAddrInfo2}))

	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}

	if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
		log.Printf("Failed to connect peerA and relay")
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
