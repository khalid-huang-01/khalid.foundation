package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

func main()  {
	//cert, err := tls.LoadX509KeyPair("./client/kubeedge/client.crt", "./client/kubeedge/client.key")
	//cert, err := tls.LoadX509KeyPair("./client/client.crt", "./client/client.key")
	cert, err := tls.LoadX509KeyPair("./client/ecdsa/client.crt", "./client/ecdsa/client.key")
	if err != nil {
		log.Println(err)
		return
	}
	caCertBytes, err := ioutil.ReadFile("./ca.crt")
	//caCertBytes, err := ioutil.ReadFile("./ca-kubeedge.crt")
	if err != nil {
		panic("unable to read client.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(caCertBytes)
	if !ok {
		panic("failed to parse root certificate")
	}

	conf := &tls.Config{
		RootCAs: clientCertPool,
		Certificates: []tls.Certificate{cert},
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
	println(string(buf[:n]))
}
