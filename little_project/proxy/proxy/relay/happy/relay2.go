package main

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"khalid.foundation/proxy/proxy/relay/utils"
)

func main() {
	listenPort := 10002
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
	)
	if err != nil {
		log.Printf("Failed to create relay server: %s", err)
		return
	}

	utils.PrintHostAddr("RELAY 2: ", host)

	select {}

}
