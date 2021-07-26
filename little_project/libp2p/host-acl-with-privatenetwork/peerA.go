package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"khalid.fondation/libp2pdemo/utils"
	"log"
)

func main() {
	key := "d6a3ab80d31ab42650da9173c764380ab7e1421b4041329ff3e1a3cbe0860f6b"
	s := ""
	//s += fmt.Sprintln("/key/swarm/psk/1.0.0/")
	//s += fmt.Sprintln("/base16/")
	s += fmt.Sprintf("%s", key)

	relayID := "QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(),
		libp2p.PrivateNetwork([]byte(s)))
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
		log.Printf("Failed to connect peerA and relay, err: %v", err)
		return
	}
	log.Println("success to connect to relay")
}