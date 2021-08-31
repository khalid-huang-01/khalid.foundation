package main

import (
	"bufio"
	"context"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"io/ioutil"
	//libp2ptls "github.com/libp2p/go-libp2p-tls"
	libp2ptlsca "khalid.fondation/libp2pdemo/go-libp2p-tls-ca"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go

func main() {
	// openssl genrsa -out rsa_private.key 2048
	certBytes, err := ioutil.ReadFile("D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\server\\server.key")
	if err != nil {
		log.Println("unable to read client.pem, error: ", err)
		return
	}
	block, _ := pem.Decode(certBytes)

	priv, err := crypto.UnmarshalECDSAPrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	os.Setenv("CERTFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\server\\server.crt")
	os.Setenv("KEYFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\server\\server.key")
	os.Setenv("CAFILE", "D:\\workspace\\gocode\\gomodule\\khalid.foundation\\little_project\\libp2p\\basic-echo-project-with-ca\\ca-ca-server.crt")

	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10001"),
		//libp2p.NoSecurity)
		//libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(libp2ptlsca.ID, libp2ptlsca.New),
	)

	if err != nil {
		panic(err)
	}

	node.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})



	// pint he node's peerInfo in multiaddr format
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address: ", addrs[0])

	// wait for a SIGINT or SIGTERM signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
	if err := node.Close(); err != nil {
		panic(err)
	}

}
// doEcho reads a line of data a stream and writes it back
func doEcho(s network.Stream) error {
	buf := bufio.NewReader(s)
	str, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	log.Printf("read: %s", str)
	_, err = s.Write([]byte(str))
	return err
}