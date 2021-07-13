package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	relayID := "QmdBMV5BHd7PV8jHJvb7pvvbzQEbjmysN4exdf3FhXgBXN"
	peerBID := "QmRVp43XvGY68BmLG56PmXF3t5n49L9sWy7Uz1NLAGww1w"
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay())
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	relayAddr := "/ip4/192.168.0.38/tcp/10001/p2p/" + relayID
	relayAddrInfo, err := utils.Addr2info(relayAddr)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
		log.Printf("Failed to connect peerA and relay")
		return
	}
	log.Println("success to connect to relay")

	// start to relay to peerB
	//relayaddrForPeerB, err := ma.NewMultiaddr("/p2p/QmcbNjSVqW6U1mAnv8koDp9hQL8K36A5f2YQqdpMUNZGXH/p2p-circuit/ipfs/QmTSUcqBkdkAmGybVK3HnqyeQk2p21QmsbQT1QPH8UfEeG")
	relayaddrForPeerB, err := ma.NewMultiaddr("/p2p/" + relayID + "/p2p-circuit/ipfs/" + peerBID)
	log.Println("relayaddr: ", relayaddrForPeerB)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	log.Println("peerBID: ", peerBID)
	id, err := peer.Decode(peerBID)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	peerBRelayInfo := peer.AddrInfo{
		ID: id,
		Addrs: []ma.Multiaddr{relayaddrForPeerB},
	}

	if err := host.Connect(context.Background(), peerBRelayInfo); err != nil {
		log.Println("Failed to connect peerA and peerB, err: ", err)
		return
	}

	s, err := host.NewStream(context.Background(), id, "/cats")
	if err != nil {
		log.Println("huh, this should have worked: ", err)
		return
	}
	s.Read(make([]byte, 1))


	select {}

}
