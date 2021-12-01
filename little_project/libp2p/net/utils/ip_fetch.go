package main

import (
	"fmt"
	"github.com/libp2p/go-netroute"
	"net"
	"strings"
)

func main()  {
	rsl := make([]net.IP, 0)

	// try to use the default ipv4/6 address
	r, err := netroute.New()
	if err != nil {
		panic(err)
	}
	_, _, localIPv4, err := r.Route(net.IPv4zero)
	if err != nil {
		panic(err)
	}

	if localIPv4.IsGlobalUnicast() {
		rsl = append(rsl, localIPv4)
	}
	fmt.Println(rsl)

	if _, _, localIPv6, err := r.Route(net.IPv6unspecified); err != nil {
		fmt.Println("1", err)
	} else if localIPv6.IsGlobalUnicast() {
		rsl = append(rsl, localIPv6)
	}

	// resolve the interface addresses
	if addrs, err := net.InterfaceAddrs(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(addrs)
		for _, addr := range addrs {
			// 127.0.0.1/8 去掉里面的/8
			ipStr := strings.Split(addr.String(), "/")[0]
			ip := net.ParseIP(ipStr)
			//fmt.Println(ip, "   ", addr.String())
			if !isIp6LinkLocal(ip) && ip.IsLoopback() {
				rsl = append(rsl, ip)
			}
		}
	}

	for _, ip := range rsl {
		fmt.Print(ip.String(), "   ")
	}
	fmt.Println("")

}

func isIp6LinkLocal(ip net.IP) bool {
	return ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast()
}

