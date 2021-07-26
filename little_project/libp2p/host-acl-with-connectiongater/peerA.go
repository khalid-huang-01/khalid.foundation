package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"io/ioutil"
	"khalid.fondation/libp2pdemo/utils"
	"log"
)

func main() {
	certBytes, err := ioutil.ReadFile("./host-acl-with-connectiongater/peerA.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}


	relayID := "QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", 10002)),
		libp2p.Identity(priv))
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	fmt.Println(host.ID())
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