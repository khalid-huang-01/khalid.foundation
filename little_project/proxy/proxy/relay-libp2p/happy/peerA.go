package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	host, err := libp2p.New(context.Background(), libp2p.EnableAutoRelay())
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	relay1Add := "/ip4/192.168.0.10/tcp/10001/p2p/QmeoSq17V9dyfryGHCcQPSPVTboeXnjK2BQ8p4caAkekz7"
	relay1AddrInfo, err := utils.Addr2info(relay1Add)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if err := host.Connect(context.Background(), *relay1AddrInfo); err != nil {
		log.Printf("Failed to connect peerA and relay1")
		return
	}
	log.Println("success to connect")
	select {}

}
