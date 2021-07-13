package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			//TLSClientConfig: &tls.Config{}, // 如果这样配置的话，会出现x509: certificate signed by unknow authority,解决方案有两种，一种是client跳过检查；另一种是将server的证书的添加到client的ca列表里
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get("https://127.0.0.1:8443")
	if err != nil {
		log.Println(err)
		return
	}

	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("%v\n", resp.Status)
	fmt.Printf(string(htmlData))
}
