package utils

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
)

func PrintHostAddr(hostName string, host host.Host) {
	for _, v := range host.Addrs() {
		fmt.Printf("%s : %v/p2p/%s\n", hostName, v, host.ID().Pretty())
	}
}

func Addr2info(addrStr string) (*peer.AddrInfo, error) {
	addr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	return peer.AddrInfoFromP2pAddr(addr)
}