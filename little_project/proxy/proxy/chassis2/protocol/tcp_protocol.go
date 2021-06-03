package protocol

import (
	"context"
	"io"
	"log"
	"net"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	handler2 "khalid.foundation/proxy/proxy/chassis2/handler"
)

type TCP struct {
	Conn net.Conn
}

// invocation里面可以通过context来传递lconn的值，然后写一个类似handler.Transport的东西，完成对数据的传输
// 这样的话，在responseCallback的可以不用做其他的事情了。但其实哪里做都可以；这套东西本来就没有面对流做支持的。
// 如果需要的话，只能在handler里做远程的数据传输，然后在tcp.callback里面做本地的数据传输，好像有点割裂
// 从语义上来看，还是把responseCallback做为这函数就可以了。
func (t *TCP) Process() {
	//inv := invocation.New(context.Background())
	inv := invocation.New(context.WithValue(context.Background(), "lconn", t.Conn))
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
		log.Printf("handle l4 proxy err: %v", data.Err)
		return data.Err
	}
	return nil
	//ep, ok := data.Result.(string)
	//if !ok {
	//	err := fmt.Errorf("result %v not string type\n", data.Result)
	//	log.Println(err)
	//	t.Conn.Close()
	//	return err
	//}
	//epSplit := strings.Split(ep, ":")
	//host := epSplit[0]
	//port, err := strconv.Atoi(epSplit[1])
	//if err != nil {
	//	err := fmt.Errorf("endpoint %s not a valid address", ep)
	//	log.Println(err)
	//	t.Conn.Close()
	//	return err
	//}
	//addr := &net.TCPAddr{
	//	IP: net.ParseIP(host),
	//	Port: port,
	//}
	//log.Printf("l4 proxy ge tserver address: %v", addr)
	//var lconn, rconn net.Conn // local Conn and remote Conn
	//lconn = t.Conn
	//for i := 0; i < 3; i++ {
	//	rconn, err = net.DialTimeout("tcp", addr.String(), time.Duration(10) * time.Second)
	//	if err == nil {
	//		break
	//	}
	//}
	//if err != nil {
	//	log.Printf("l4 proxy dia server error: %v", err)
	//	t.Conn.Close()
	//	return err
	//}
	//// 这里不作作阻塞吗？
	//go t.pipe(lconn, rconn)
	//go t.pipe(rconn, lconn)
	//return nil
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
