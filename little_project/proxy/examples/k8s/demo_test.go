package k8s

import (
	"testing"

	utilnet "k8s.io/apimachinery/pkg/util/net"
)

func TestDemo(t *testing.T) {
	advertiseAddress, _ := utilnet.ChooseHostInterface()
	t.Log(advertiseAddress)
}
