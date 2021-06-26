package main

import (
	"log"

	"khalid.foundation/proxy/proxy/ha-relay-server/activestandby"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	host, err := activestandby.NewServer()
	if err != nil {
		panic(err)
	}

	utils.PrintHostAddr("RELAY : ", host)
	log.Println("ID: ", host.ID())

	select {}
}