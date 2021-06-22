package main

import (
	"context"
	"encoding/pem"
	"io/ioutil"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"khalid.foundation/proxy/proxy/relay-libp2p/utils"
)

func main() {
	// 证书加载
	certBytes, err := ioutil.ReadFile("./client.key")
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

	relayID := "QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	// libp2p.ListenAddrs的作用是什么 => 启动服务，这样别人才能通过stream连接自己，这个是默认启动的，不用配置
	host, err := libp2p.New(context.Background(), libp2p.EnableRelay(),
		libp2p.Ping(true),
		libp2p.Identity(priv))
	if err != nil {
		log.Printf("Failed to create h1: %s", err)
		return
	}
	relayAddr := "/ip4/127.0.0.1/tcp/10001/p2p/" + relayID
	relayAddrInfo, err := utils.Addr2info(relayAddr)
	if err != nil {
		log.Println("err: ", err)
		return
	}

	if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
		log.Println("Failed to connect relay, error: ", err)
		return
	}

	//ticker := time.NewTicker(time.Second * 3)
	//go func() {
	//	for {
	//		select {
	//		case <-ticker.C:
	//			if err := host.Connect(context.Background(), *relayAddrInfo); err != nil {
	//				log.Println("Failed to connect relay, error: ", err)
	//			} else {
	//				log.Println("success connected")
	//			}
	//		default:
	//			time.Sleep(time.Second * 2)
	//		}
	//	}
	//}()


	host.SetStreamHandler("/cats", func(s network.Stream) {
		log.Println("Meow! It worked!")
		s.Close()
	})

	utils.PrintHostAddr("Server peerB : ", host)
	log.Println("ID: ", host.ID())
	select {}
}
