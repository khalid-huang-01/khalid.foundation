package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"khalid.fondation/libp2pdemo/utils"
	"log"
)

func main() {

	relayID := "QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	// libp2p.ListenAddrs的作用是什么 => 启动服务，这样别人才能通过stream连接自己，这个是默认启动的，不用配置
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay())
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	relayAddr := "/ip4/192.168.0.10/tcp/10001/p2p/" + relayID
	relayAddrInfo, err := utils.Addr2info(relayAddr)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
		log.Println("Failed to connect relay, error: ", err)
		return
	}
	log.Println("success to connect to relay")
}