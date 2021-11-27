package main

import (
	"context"
	"fmt"
	"khalid.fondation/libp2pdemo/utils"
	"log"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	listenPort := 10001
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
		libp2p.Ping(true),
	)
	if err != nil {
		log.Printf("Failed to create relay-libp2p server: %s", err)
		return
	}

	utils.PrintHostAddr("RELAY : ", host)
	log.Println("ID: ", host.ID())

	select {}

}
