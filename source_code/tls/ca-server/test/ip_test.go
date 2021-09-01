package test

import (
	"fmt"
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
		fmt.Println(ipnet.IP.String())
	}
}
