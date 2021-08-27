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

type Gater struct {
}

// InterceptPeerDial 这个是做outbound 流量的控制的,当这个返回false的时候,节点之间的connection是无法建立的,可以通过
// netstat -anp | grep 10001 查看验证
// 但是remote peerB connect不会返回失败
func (g *Gater) InterceptPeerDial(p peer.ID) (allow bool) {
	fmt.Println("remote peer ID: ", p)
	if p.String() == "QmbyWhnkqUQArfzukTRTQqnzuAVQfSTuijTN1iHsK5CU8z" {
		return true
	}
	return false
}

func (g *Gater) InterceptAddrDial(p peer.ID,_ ma.Multiaddr) (allow bool) {
	return g.InterceptPeerDial(p)
}

// InterceptAccept 这个是做inbound流量的,但这个false, 节点之间的connection也是无法建立的,而且remote Peer connect
// 会返回失败
func (g *Gater) InterceptAccept(connAddr network.ConnMultiaddrs) (allow bool) {
	fmt.Printf("%v\n", connAddr.RemoteMultiaddr())
	//return false
	return true
}


func (g *Gater) InterceptSecured(network.Direction, peer.ID, network.ConnMultiaddrs) (allow bool) {
	return true
}

func (g *Gater) InterceptUpgraded(n network.Conn) (allow bool, reason control.DisconnectReason) {
	return true, 0
}
