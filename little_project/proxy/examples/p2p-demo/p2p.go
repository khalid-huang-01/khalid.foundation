package main

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
)

func main() {
	certBytes, err := ioutil.ReadFile("./examples/p2p-demo/server.key")
	if err != nil {
		log.Println("unable to read client.pem")
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalRsaPrivateKey(block.Bytes)
	//priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("unable to unmarshal, err: ", err)
		return
	}
	log.Println("priv: ", priv)
	h1, err := libp2p.New(context.Background(), libp2p.Identity(priv))
	if err != nil {
		log.Println("faile dto create h1, err:", err)
		return
	}
	log.Println(h1.ID().Pretty())
}
