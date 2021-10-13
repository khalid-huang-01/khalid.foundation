package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	usages := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	fmt.Println(usages)
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/hello", nil)
	if err != nil {
		panic(err)
	}
	//usagesStr, err := json.Marshal(&usages)
	//if err != nil {
	//	panic(err)
	//}
	//req.Header.Set("ExtKeyUsages", string(usagesStr))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	//response, _ := http.Get("http://127.0.0.1:8080/hello")

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
