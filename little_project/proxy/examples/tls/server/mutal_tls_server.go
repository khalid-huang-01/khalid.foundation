package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
)

// 一般客户端本地是有CA的
func main() {
	cert, err := tls.LoadX509KeyPair("./examples/tls/server/server.pem", "./examples/tls/server/server.key")
	if err != nil {
		log.Println(err)
		return
	}

	certBytes, err := ioutil.ReadFile("./examples/tls/client/client.pem")
	if err != nil {
		panic("unable to read client.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: clientCertPool,
	}

	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println("err")
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn1(conn)
	}
}

func handleConn1(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}