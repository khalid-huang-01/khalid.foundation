package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	http.HandleFunc("/hello", hello)
	log.Println("http server listen on 30002")
	err := http.ListenAndServe(":30002", nil)
	if err != nil {
		log.Fatal("ListenAndServer fail: ", err)
	}
}
