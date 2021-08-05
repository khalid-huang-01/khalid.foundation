package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/pnet"
	"khalid.fondation/libp2pdemo/utils"
	"log"
)

func main() {
	key := "2bc218803c6d2a57e1709a1199997c5cee6414201d9c1d13488d12cfd61cbd96"
	s := ""
	s += fmt.Sprintln("/key/swarm/psk/1.0.0/")
	s += fmt.Sprintln("/base16/")
	s += fmt.Sprintf("%s", key)
	psk, err := pnet.DecodeV1PSK(bytes.NewBuffer([]byte(s)))
	if err != nil {
		panic(err)
	}

	relayID := "QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(),
		libp2p.PrivateNetwork(psk))
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