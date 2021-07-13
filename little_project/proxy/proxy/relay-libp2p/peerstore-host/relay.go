package main

import (
	"context"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"time"

	ds "github.com/ipfs/go-datastore"
	badger "github.com/ipfs/go-ds-badger"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoreds"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"

	"log"
)

func main() {
	// 证书加载
	certBytes, err := ioutil.ReadFile("./server-1.key")
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

	listenPort := 10001
	ps := NewPeerstoreWithBadger()
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
		libp2p.Peerstore(ps),
		libp2p.Identity(priv),
		)
	if err != nil {
		log.Printf("Failed to create relay-libp2p server: %s", err)
		return
	}

	utils.PrintHostAddr("RELAY : ", host)
	log.Println("ID: ", host.ID())

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Println(ps.Peers())
				log.Println(ps.Addrs("QmU2Nq2RJSC4os9WGwN5J5fbxTxsND26NmUTds1JSzWzuE"))
				log.Printf("%+v", ps.PeerInfo("QmU2Nq2RJSC4os9WGwN5J5fbxTxsND26NmUTds1JSzWzuE"))
			default:
				time.Sleep(time.Second * 2)
			}
		}
	}()
	select {}

}

func badgerStore() (ds.Batching, func()) {
	dataPath := "D:\\workspace\\data\\badger"
	store, err := badger.NewDatastore(dataPath, &badger.DefaultOptions)
	if err != nil {
		panic(err)
	}
	closer := func() {
		store.Close()
	}
	return store, closer
}

func NewPeerstoreWithBadger() peerstore.Peerstore {
	store, _ := badgerStore()
	ps, err := pstoreds.NewPeerstore(context.Background(), store, pstoreds.DefaultOpts())
	if err != nil {
		panic(err)
	}
	return ps
}