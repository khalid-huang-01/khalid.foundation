package activestandby

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	circuit "github.com/libp2p/go-libp2p-circuit"
	"github.com/libp2p/go-libp2p-core/host"
)

// TODO 要添加选主和readyz接口
func NewServer() (host.Host, error) {
	// TODO 选主
	rea

	listenPort := 10001
	host, err := libp2p.New(
		context.Background(),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.EnableRelay(circuit.OptHop),
		libp2p.ForceReachabilityPrivate(),
	)
	if err != nil {
		log.Printf("Failed to create relay-libp2p server: %s", err)
		return nil, err
	}

	return host, nil
}

