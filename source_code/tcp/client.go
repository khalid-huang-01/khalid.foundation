package main

import (
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.0.10:30002")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// read echo service
	buf := make([]byte, 100)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	println(string(buf[:n]))

	// write hello
	n, err = conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
	}

	// read  hello
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	println(string(buf[:n]))
}