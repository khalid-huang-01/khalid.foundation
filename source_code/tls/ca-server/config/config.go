package config

import (
	"encoding/pem"
	"io/ioutil"
	"k8s.io/klog/v2"
)

var Config Configure

type Configure struct {
	TLSCAFile string
	TLSCAKeyFile string
	TLSCertFile string
	TLSPrivateKeyFile string

	Ca []byte
	CaKey []byte
	Cert []byte
	Key []byte

	// server tls
	DNSNames []string
	AdvertiseAddress []string
}

func InitConfigure() {
	Config = Configure{
		TLSCAFile: "D:\\workspace\\gocode\\gomodule\\data\\ca\\rootCA.crt",
		TLSCAKeyFile: "D:\\workspace\\gocode\\gomodule\\data\\ca\\rootCA.key",
		TLSCertFile: "D:\\workspace\\gocode\\gomodule\\data\\ca\\server.crt",
		TLSPrivateKeyFile: "D:\\workspace\\gocode\\gomodule\\data\\ca\\server.key",
		DNSNames: []string{"localhost"},
		AdvertiseAddress: []string{"127.0.0.1"},
	}

	ca, err := ioutil.ReadFile(Config.TLSCAFile)
	if err == nil {
		block, _ := pem.Decode(ca)
		ca = block.Bytes
		klog.Info("Succeed in loading CA certificate from local directory")
	}

	caKey, err := ioutil.ReadFile(Config.TLSCAKeyFile)
	if err == nil {
		block, _ := pem.Decode(caKey)
		caKey = block.Bytes
		klog.Info("Succeed in loading CA key from local directory")
	}

	if ca != nil && caKey != nil {
		Config.Ca = ca
		Config.CaKey = caKey
	} else if !(ca == nil && caKey == nil) {
		klog.Fatal("Both of ca and caKey should be specified!")
	}

	cert, err := ioutil.ReadFile(Config.TLSCertFile)
	if err == nil {
		block, _ := pem.Decode(cert)
		cert = block.Bytes
		klog.Info("Succeed in loading certificate from local directory")
	}
	key, err := ioutil.ReadFile(Config.TLSPrivateKeyFile)
	if err == nil {
		block, _ := pem.Decode(key)
		key = block.Bytes
		klog.Info("Succeed in loading private key from local directory")
	}

	if cert != nil && key != nil {
		Config.Cert = cert
		Config.Key = key
	} else if !(cert == nil && key == nil) {
		klog.Fatal("Both of cert and key should be specified!")
	}
}

func UpdateConfig(ca, caKey, cert, key []byte) {
	if ca != nil {
		Config.Ca = ca
	}
	if caKey != nil {
		Config.CaKey = caKey
	}
	if cert != nil {
		Config.Cert = cert
	}
	if key != nil {
		Config.Key = key
	}
}
