package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
)

func main()  {
	address := "127.0.0.1:8080"
	http.HandleFunc("/hello", Hello)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		panic(err)
	}
}

func Hello(w http.ResponseWriter, req *http.Request) {
	usagesStr := req.Header.Get("ExtKeyUsages")
	if usagesStr == "" {
		fmt.Println(true)
	}
	fmt.Println(usagesStr)
	var usages []x509.ExtKeyUsage
	err := json.Unmarshal([]byte(usagesStr), &usages)
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Println(usages)
	w.Write([]byte("hello world"))
}
