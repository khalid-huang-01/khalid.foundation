package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p-core/control"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"khalid.fondation/libp2pdemo/utils"
	"log"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	// openssl genrsa -out rsa_private.key 2048
	certBytes, err := ioutil.ReadFile("./host-acl-with-connectiongater/server.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	listenPort := 10001

	host, err := libp2p.New(
		context.Background(),
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
		libp2p.ConnectionGater(&Gater{}),
	)
	if err != nil {
		log.Printf("Failed to create relay-libp2p server: %s", err)
		return
	}

	utils.PrintHostAddr("RELAY : ", host)
	log.Println("ID: ", host.ID())

	select {}

}

type Gater struct {
}

func (g *Gater) InterceptPeerDial(p peer.ID) (allow bool) {
	return true
}

func (g *Gater) InterceptAddrDial(peer.ID, ma.Multiaddr) (allow bool) {
	fmt.Println("addr dial")
	return true
}

func (g *Gater) InterceptAccept(connAddr network.ConnMultiaddrs) (allow bool) {
	if connAddr.RemoteMultiaddr().String() == "/ip4/192.168.0.10/tcp/10002" {
		return true
	}
	return false
}


func (g *Gater) InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool) {
	return true
}

func (g *Gater) InterceptUpgraded(n network.Conn) (allow bool, reason control.DisconnectReason) {
	fmt.Println("n", n)
	return true, 0
}
