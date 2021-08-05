package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/pnet"
	"io/ioutil"
	"khalid.fondation/libp2pdemo/utils"
	"log"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

func main() {
	certBytes, err := ioutil.ReadFile("./host-acl-with-privatenetwork/server.key")
	if err != nil {
		log.Println("unable to read client.pem")
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	listenPort := 10001

	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		panic(err)
	}
	s := ""
	s += fmt.Sprintln("/key/swarm/psk/1.0.0/")
	s += fmt.Sprintln("/base16/")
	s += fmt.Sprintf("%s", hex.EncodeToString(key))
	fmt.Println(s)
	psk, err := pnet.DecodeV1PSK(bytes.NewBuffer([]byte(s)))
	if err != nil {
		panic(err)
	}


	host, err := libp2p.New(
		context.Background(),
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
		libp2p.PrivateNetwork(psk),
	)
	if err != nil {
		log.Printf("Failed to create relay-libp2p server: %s", err)
		return
	}

	utils.PrintHostAddr("RELAY : ", host)
	log.Println("ID: ", host.ID())

	select {}

}
