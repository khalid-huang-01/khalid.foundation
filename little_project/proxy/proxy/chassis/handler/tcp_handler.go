package handler

import (
	"log"

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
	r := &invocation.Response{
		Result: i.Endpoint,
	}

	if err := cb(r); err != nil {
		log.Println("failed to cb: ", err)
	}

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
