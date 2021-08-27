package main

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	libp2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/multiformats/go-multiaddr"
	"io/ioutil"
	"log"
)

func main()  {
	ctx := context.Background()

	// 这里使用libp2p.NoSecurity，如果对端不使用这个或者Security里面的协议不一样的化，两者就无法建立连接
	// 因为没有对应的协议做Secrutiy
	node, err := libp2p.New(ctx, libp2p.Ping(false),
		//libp2p.NoSecurity)
		libp2p.Security(libp2ptls.ID, libp2ptls.New))
	if err != nil {
		panic(err)
	}

	//if len(os.Args) <= 1 {
	//	panic("Please provide the peer addr")
	//}
	//addr, err := multiaddr.NewMultiaddr(os.Args[1])
	addrStr := "/ip4/192.168.0.38/tcp/10001/p2p/QmbSUTgoPDgRqP5S1Zz2fJJhtg8MFiQna3XAQTQRk9nDSG"
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	peer, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *peer); err != nil {
		panic(err)
	}

	s, err := node.NewStream(ctx, peer.ID, "/echo/1.0.0")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Sender syaning hello")
	_, err = s.Write([]byte("Hello, world\n"))

	if err != nil {
		log.Println(err)
		return
	}

	out, err := ioutil.ReadAll(s)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("read reply: %q\n", out)

	if err := node.Close(); err != nil {
		panic(err)
	}}
