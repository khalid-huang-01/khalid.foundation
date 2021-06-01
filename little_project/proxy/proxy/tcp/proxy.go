package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
)

// laddr本地监听地址，30001 //local
// raddr目标访问地址：30002 //remote
// lconn, client与本服务tcp server建立的连接
// rconn, 本服务与目标tcp 服务建立的连接
type Proxy struct {
	laddr, raddr string
	lconn, rconn net.Conn // 这里也可以是io.ReadWriteClose
	errCh chan string
}

func New(lconn net.Conn, laddr, raddr string) *Proxy {
	return &Proxy{
		lconn: lconn,
		laddr: laddr,
		raddr: raddr,
		errCh: make(chan string),
	}
}

// 本地监听30001，代理到本地30002
func (p *Proxy) Start() {
	defer p.lconn.Close()
	var err error
	p.rconn, err = net.Dial("tcp", p.raddr)
	if err != nil {
		log.Println("Remote connection failed: ", err)
		return
	}
	defer p.rconn.Close()
	go p.pipe(p.lconn, p.rconn) // 把本地接收到的请求代理到远程服务
	go p.pipe(p.rconn, p.lconn) // 把从远程服务接收到的响应写给本地连接

	// 这里粗处理，当流里面读取到EOF，也就是一个连接正常结束的时候，也是来到这里的
	errStr := <-p.errCh
	log.Println("proxy exit err: ", errStr)
}

func (p *Proxy) pipe(src, des io.ReadWriter) {
	islocal := src == p.lconn // 判断是否是本地发送远程的，反之就是回回来的响应
	var dataDirection string
	if islocal {
		dataDirection = ">>> %d bytes send\n"
	} else {
		dataDirection = "<<< %d bytes received\n"
	}
	buff := make([]byte, 0xffff) //64K
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.errCh <- fmt.Sprintf("Read failed, err: %s", err)
			return
		}
		b := buff[:n]
		log.Printf(dataDirection, n)
		n, err = des.Write(b)
		if err != nil {
			p.errCh <- fmt.Sprintf("Write failed, err: %s", err)
			return
		}
	}

}