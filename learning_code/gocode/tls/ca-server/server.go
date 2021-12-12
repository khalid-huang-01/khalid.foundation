package main

import (
	"k8s.io/klog/v2"
	"khalid.jobs/caserver/config"
	"khalid.jobs/caserver/httpserver"
)

func main() {
	config.InitConfigure()

	if err := httpserver.PrepareAllCerts(); err != nil {
		klog.Fatal(err)
	}

	if err := httpserver.GenerateToken(); err != nil {
		klog.Fatal(err)
	}

	httpserver.StartHTTPServer()
}