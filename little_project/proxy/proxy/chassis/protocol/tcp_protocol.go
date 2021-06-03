package protocol

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	handler2 "khalid.foundation/proxy/proxy/chassis/handler"
)

type TCP struct {
	Conn net.Conn
}

func (t *TCP) Process() {
	inv := invocation.New(context.Background())
	inv.MicroServiceName = "edge.default.svc.cluster.local:12345"
	inv.SourceServiceID = ""
	inv.Strategy = loadbalancer.StrategyRoundRobin
	// 这里是tcp，所以使用自己定义的handler2.L4ProxyHandlerName， 在http里面，使用默认的handler.Transport
	c, err := handler.CreateChain(common.Consumer, "tcp", handler.Loadbalance, handler2.L4ProxyHandlerName)
	if err != nil {
		log.Printf("create handler chaiin error: %v \n", err)
		t.Conn.Close()
	}
	c.Next(inv, t.responseCallback)
}

func (t *TCP) responseCallback(data *invocation.Response) error {
	if data.Err != nil {
		log.Println("handler l4 proxy err: ", data.Err)
		t.Conn.Close()
		return data.Err
	}

	ep, ok := data.Result.(string)
	if !ok {
		err := fmt.Errorf("result %v not string type\n", data.Result)
		log.Println(err)
		t.Conn.Close()
		return err
	}
	epSplit := strings.Split(ep, ":")
	host := epSplit[0]
	port, err := strconv.Atoi(epSplit[1])
	if err != nil {
		err := fmt.Errorf("endpoint %s not a valid address", ep)
		log.Println(err)
		t.Conn.Close()
		return err
	}
	addr := &net.TCPAddr{
		IP: net.ParseIP(host),
		Port: port,
	}
	log.Printf("l4 proxy ge tserver address: %v", addr)
	var lconn, rconn net.Conn // local Conn and remote Conn
	lconn = t.Conn
	for i := 0; i < 3; i++ {
		rconn, err = net.DialTimeout("tcp", addr.String(), time.Duration(10) * time.Second)
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Printf("l4 proxy dia server error: %v", err)
		t.Conn.Close()
		return err
	}
	// 这里不作作阻塞吗？
	go t.pipe(lconn, rconn)
	go t.pipe(rconn, lconn)
	return nil
}

func (t *TCP) pipe(src, des io.ReadWriteCloser) {
	//_, err := io.Copy(des, src)
	//if err != nil {
	//	fmt.Println("read error: ", err )
	//}

	buff := make([]byte, 0xffff) //64K
	for {
		n, err := src.Read(buff)
		if err != nil {
			if err != io.EOF {
				log.Println("read error: ", err)
			}
			src.Close()
			des.Close()
			break
		}
		_, err = des.Write(buff[:n])
		if err != nil {
			log.Println("write error: ", err)
			src.Close()
			des.Close()
			break
		}
	}
}
