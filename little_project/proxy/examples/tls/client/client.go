package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
)

func main()  {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conf.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		fmt.Println(rawCerts)
		return fmt.Errorf("verify faile")
	}

	conn, err := tls.Dial("tcp", "127.0.0.1:443", conf)
	if err != nil {
		log.Println("Dial faile: ", err)
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
