package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main()  {
	//cert, err := tls.LoadX509KeyPair("./client/ca-server/client.crt", "./client/ca-server/client.key")
	cert, err := tls.LoadX509KeyPair("./client/kubeedge/client.crt", "./client/kubeedge/client.key")
	//cert, err := tls.LoadX509KeyPair("./client/client.crt", "./client/client.key")
	//cert, err := tls.LoadX509KeyPair("./client/tmp/server.crt", "./client/tmp/server.key")
	if err != nil {
		log.Println(err)
		return
	}
	//caCertBytes, err := ioutil.ReadFile("./ca-ca-server.crt")
	//caCertBytes, err := ioutil.ReadFile("./client/tmp/ca.crt")
	caCertBytes, err := ioutil.ReadFile("./client/tmp/rootCA.crt")
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
		InsecureSkipVerify: true, // Not actually skipping, we check the cert in VerifyPeerCertificate
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			// Code copy/pasted and adapted from
			// https://github.com/golang/go/blob/81555cb4f3521b53f9de4ce15f64b77cc9df61b9/src/crypto/tls/handshake_client.go#L327-L344, but adapted to skip the hostname verification.
			// See https://github.com/golang/go/issues/21971#issuecomment-412836078.

			// If this is the first handshake on a connection, process and
			// (optionally) verify the server's certificates.
			certs := make([]*x509.Certificate, len(rawCerts))
			for i, asn1Data := range rawCerts {
				cert, err := x509.ParseCertificate(asn1Data)
				if err != nil {
					return fmt.Errorf("bitbox/electrum: failed to parse certificate from server: " + err.Error())
				}
				certs[i] = cert
			}

			fmt.Println("verify------------------------")

			opts := x509.VerifyOptions{
				Roots:         clientCertPool,
				CurrentTime:   time.Now(),
				DNSName:       "", // <- skip hostname verification
				Intermediates: x509.NewCertPool(),
			}

			for i, cert := range certs {
				if i == 0 {
					continue
				}
				opts.Intermediates.AddCert(cert)
			}
			_, err := certs[0].Verify(opts)
			return err
		},
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
