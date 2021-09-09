package test

import (
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"testing"
)

func TestIP(t *testing.T) {
	advertiseAddress, _ := utilnet.ChooseHostInterface()
	t.Log(advertiseAddress)
	//addrs, err := net.InterfaceAddrs()
	//if err != nil {
	//	panic(err)
	//}
	//for _, address := range addrs {
	//	ipnet := address.(*net.IPNet)
	//	fmt.Println(ipnet.IP.String())
	//}
}
