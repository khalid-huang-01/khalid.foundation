package test

import (
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	//advertiseAddress, _ := utilnet.ChooseHostInterface()
	//t.Log(advertiseAddress)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, address := range addrs {
		ipnet := address.(*net.IPNet)

		t.Log(ipnet.IP.String())
	}
}
