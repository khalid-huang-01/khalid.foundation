package handler

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
)

const L4ProxyHandlerName = "l4Proxy"

type L4ProxyHandler struct {
}

func (l *L4ProxyHandler) Name() string {
	return L4ProxyHandlerName
}

func (l *L4ProxyHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	lconn := i.Ctx.Value("lconn").(net.Conn)

	epSplit := strings.Split(i.Endpoint, ":")
	host := epSplit[0]
	port, err := strconv.Atoi(epSplit[1])
	if err != nil {
		err := fmt.Errorf("endpoint %s not a valid address", i.Endpoint)
		log.Println(err)
		lconn.Close()
	}
	addr := &net.TCPAddr{
		IP: net.ParseIP(host),
		Port: port,
	}
	log.Printf("l4 proxy ge tserver address: %v", addr)
	var  rconn net.Conn // local Conn and remote Conn
	for i := 0; i < 3; i++ {
		rconn, err = net.DialTimeout("tcp", addr.String(), time.Duration(10) * time.Second)
		if err == nil {
			break
		}
	}

	r := &invocation.Response{}

	if err != nil {
		log.Printf("l4 proxy dia server error: %v", err)
		lconn.Close()
		r.Err = err
		cb(r)
		return
	}

	// 这里要协调好，只有两个都结束了，通过传递信息才调用cb
	// 当前先简单处理
	go pipe(lconn, rconn)
	go pipe(rconn, lconn)

	cb(r)
}

// 这里要了解下读或者写结束的时候，会返回什么结束码
func pipe(src, des io.ReadWriteCloser) {
//	// 从响应上来看，在STOP的时候，没有把数据回传 原因在于说：io.Copy结束的时候err是为空的
	// 如果中途断了，可能需要做一些重新连接的操作？
	_, err := io.Copy(des, src)
	if err != nil {
		fmt.Println("read error: ", err )
	}
	src.Close()
	des.Close()

	//buff := make([]byte, 0xffff) //64K
	//for {
	//	n, err := src.Read(buff)
	//	if err != nil {
	//		if err != io.EOF {
	//			log.Println("read error: ", err)
	//		}
	//		src.Close()
	//		des.Close()
	//		break
	//	}
	//	_, err = des.Write(buff[:n])
	//	if err != nil {
	//		log.Println("write error: ", err)
	//		src.Close()
	//		des.Close()
	//		break
	//	}
	//}
}


func newL4ProxyHandler() handler.Handler {
	return &L4ProxyHandler{}
}

func init() {
	log.Println("register l4 proxy handler")
	err := handler.RegisterHandler(L4ProxyHandlerName, newL4ProxyHandler)
	if err != nil {
		log.Println("register l4 proxy handler err: ", err)
		panic("failed register")
	}
}
