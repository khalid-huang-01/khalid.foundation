package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// https://docs.libp2p.io/tutorials/getting-started/go

func main() {
	ctx := context.Background()
	// 为了更好的知道ping服务，这里禁止了直接使用内置的ping protocol
	node, err := libp2p.New(ctx, libp2p.Ping(false),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/10002"),
		libp2p.EnableRelay(circuit.OptActive),
		libp2p.ForceReachabilityPrivate())

	if err != nil {
		panic(err)
	}
	fmt.Println(node.ID().Pretty())

	node.SetStreamHandler("/echo/1.0.0", func(s network.Stream) {
		log.Println("listener received new stream")
		if err := doEcho(s); err != nil {
			log.Println(err)
			s.Reset()
		} else {
			s.Close()
		}
	})
	// server与relay建立连接，方便其他client通过server访问进来（这里的server和client其实不是我们常说的那种，在libp2里面的都是agent，只是agent的功能不大一样）
	//relayStr := "/ip4/192.168.0.10/tcp/10001/p2p/Qma57EzsaP5FyNpqduXpCgR9xUkWoXkkKoZurTH9rqPBgC"
    relayStr := "/ip4/119.13.84.169/tcp/10001/p2p/QmUJs7ut4U376MxopVf3qTxVAK1AcrLUHDMnr3nJobTsqU"
	relayAddr, err := multiaddr.NewMultiaddr(relayStr)
	if err != nil {
		fmt.Println("1")
		panic(err)
	}
	relayPeer, err := peer.AddrInfoFromP2pAddr(relayAddr)
	if err != nil {
		fmt.Println("2")
		panic(err)
	}
	if len(node.Network().ConnsToPeer(relayPeer.ID)) == 0 {
		err := node.Connect(context.Background(), *relayPeer)
		if err != nil {
			fmt.Println("3")
			panic(err)
		}
		fmt.Println(InfoFromHostAndRelay(node, relayPeer))
	}



	// pint he node's peerInfo in multiaddr format
	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		fmt.Println("libp2p node address: ", addr)
	}
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

func InfoFromHostAndRelay(host host.Host, relay *peer.AddrInfo) *peer.AddrInfo {
	p2pProto := multiaddr.ProtocolWithCode(multiaddr.P_P2P)
	circuitProto := multiaddr.ProtocolWithCode(multiaddr.P_CIRCUIT)
	peerAddrInfo := &peer.AddrInfo{
		ID:    host.ID(),
		Addrs: host.Addrs(),
	}
	for _, v := range relay.Addrs {
		circuitAddr, err := multiaddr.NewMultiaddr(v.String() + "/" + p2pProto.Name + "/" + relay.ID.String() + "/" + circuitProto.Name)
		if err != nil {
			fmt.Printf("New multi addr err: %v\n", err)
			continue
		}
		peerAddrInfo.Addrs = append(peerAddrInfo.Addrs, circuitAddr)
	}
	return peerAddrInfo
}
