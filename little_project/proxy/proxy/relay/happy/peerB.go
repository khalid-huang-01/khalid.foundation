package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
)

func main() {
	host, err := libp2p.New(context.Background(), libp2p.EnableAutoRelay())
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}

}
