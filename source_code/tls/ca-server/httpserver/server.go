/*
Copyright 2020 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httpserver

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"khalid.jobs/caserver/config"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
)

const (
	certificateBlockType = "CERTIFICATE"
)

// StartHTTPServer starts the http service
func StartHTTPServer() {
	router := mux.NewRouter()
	router.HandleFunc("/edge.crt", edgeCoreClientCert).Methods("GET")
	router.HandleFunc("/ca.crt", getCA).Methods("GET")

	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 10002)

	cert, err := tls.X509KeyPair(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: config.Config.Cert}), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: config.Config.Key}))

	if err != nil {
		klog.Fatal(err)
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequestClientCert,
		},
	}
	klog.Fatal(server.ListenAndServeTLS("", ""))
}

// getCA returns the caCertDER
func getCA(w http.ResponseWriter, r *http.Request) {
	caCertDER := config.Config.Ca
	if _, err := w.Write(caCertDER); err != nil {
		klog.Errorf("failed to write caCertDER, err: %v", err)
	}
}

// edgeCoreClientCert will verify the certificate of EdgeCore or token then create EdgeCoreCert and return it
func edgeCoreClientCert(w http.ResponseWriter, r *http.Request) {
	if cert := r.TLS.PeerCertificates; len(cert) > 0 {
		if err := verifyCert(cert[0]); err != nil {
			klog.Errorf("failed to sign the certificate for edgenode: %s, failed to verify the certificate", "edge-node")
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				klog.Errorf("failed to write response, err: %v", err)
			}
		} else {
			signEdgeCert(w, r)
		}
		return
	}
	if verifyAuthorization(w, r) {
		signEdgeCert(w, r)
	} else {
		klog.Errorf("failed to sign the certificate for edgenode: %s, invalid token", "edge-node")
	}
}

// verifyCert verifies the edge certificate by CA certificate when edge certificates rotate.
func verifyCert(cert *x509.Certificate) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(pem.EncodeToMemory(&pem.Block{Type: certificateBlockType, Bytes: config.Config.Ca}))
	if !ok {
		return fmt.Errorf("failed to parse root certificate")
	}
	opts := x509.VerifyOptions{
		Roots:     roots,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	if _, err := cert.Verify(opts); err != nil {
		return fmt.Errorf("failed to verify edge certificate: %v", err)
	}
	return nil
}

// verifyAuthorization verifies the token from EdgeCore CSR
func verifyAuthorization(w http.ResponseWriter, r *http.Request) bool {
	authorizationHeader := r.Header.Get("authorization")
	if authorizationHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Errorf("failed to write http response, err: %v", err)
		}
		return false
	}
	bearerToken := strings.Split(authorizationHeader, " ")
	if len(bearerToken) != 2 {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Errorf("failed to write http response, err: %v", err)
		}
		return false
	}
	token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		caKey := config.Config.CaKey
		return caKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
				klog.Errorf("Wrire body error %v", err)
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Errorf("Wrire body error %v", err)
		}

		return false
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Invalid authorization token")); err != nil {
			klog.Errorf("Wrire body error %v", err)
		}
		return false
	}
	return true
}

// signEdgeCert signs the CSR from EdgeCore
func signEdgeCert(w http.ResponseWriter, r *http.Request) {
	csrContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		klog.Errorf("fail to read file when signing the cert for edgenode:%s! error:%v", "edge-node", err)
	}
	csr, err := x509.ParseCertificateRequest(csrContent)
	if err != nil {
		klog.Errorf("fail to ParseCertificateRequest of edgenode: %s! error:%v","edge-node", err)
	}
	clientCertDER, err := signCerts(csr, csr.PublicKey)
	if err != nil {
		klog.Errorf("fail to signCerts for edgenode:%s! error:%v", "edge-node", err)
	}

	if _, err := w.Write(clientCertDER); err != nil {
		klog.Errorf("wrire error %v", err)
	}
}

// signCerts will create a certificate for EdgeCore
func signCerts(csr *x509.CertificateRequest, pbKey crypto.PublicKey) ([]byte, error) {
	cfgs := &certutil.Config{
		CommonName:   csr.Subject.CommonName,
		Organization: csr.Subject.Organization,
		//Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		Usages:       []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
		AltNames: certutil.AltNames{
			DNSNames: csr.DNSNames,
			IPs:      csr.IPAddresses,
		},
	}
	clientKey := pbKey

	ca := config.Config.Ca
	caCert, err := x509.ParseCertificate(ca)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseCertificate: %v", err)
	}

	caKeyDER := config.Config.CaKey
	caKey, err := x509.ParseECPrivateKey(caKeyDER)
	if err != nil {
		return nil, fmt.Errorf("unable to ParseECPrivateKey: %v", err)
	}

	var edgeCertSigningDuration time.Duration = 365 * 24
	certDER, err := NewCertFromCa(cfgs, caCert, clientKey, caKey, edgeCertSigningDuration) //crypto.Signer(caKey)
	if err != nil {
		return nil, fmt.Errorf("unable to NewCertFromCa: %v", err)
	}

	return certDER, err
}

// PrepareAllCerts check whether the certificates exist in the local directory,
// and then check whether certificates exist in the secret, generate if they don't exist
func PrepareAllCerts() error {
	// Check whether the ca exists in the local directory
	if config.Config.Ca == nil && config.Config.CaKey == nil {
		klog.Info("Ca and CaKey don't exist in local directory, and will create new")
		// Check whether the ca exists in the secret
		caDER, caKey, err := NewCertificateAuthorityDer()
		if err != nil {
			klog.Errorf("failed to create Certificate Authority, error: %v", err)
			return err
		}

		caKeyDER, err := x509.MarshalECPrivateKey(caKey.(*ecdsa.PrivateKey))
		if err != nil {
			klog.Errorf("failed to convert an EC private key to SEC 1, ASN.1 DER form, error: %v", err)
			return err
		}
		caCert, err  := x509.ParseCertificate(caDER)
		if err != nil {
			klog.Errorf("failed to parse Certificate %v", err)
			return err
		}
		config.UpdateConfig(caDER, caKeyDER, nil, nil)
		if err = WriteKeyAndCert(config.Config.TLSCAKeyFile, config.Config.TLSCAFile, caKey, caCert); err != nil {
			return fmt.Errorf("failed to save the edge key and certificate to file: %s, error: %v", config.Config, err)
		}
	}

	// Check whether the CloudCore certificates exist in the local directory
	if config.Config.Key == nil && config.Config.Cert == nil {
		klog.Infof("CloudCoreCert and key don't exist in local directory,and will be signed by CA")
		certDER, keyDER, err := SignCerts()
		if err != nil {
			klog.Errorf("failed to sign a certificate, error: %v", err)
			return err
		}
		key, err :=  x509.ParseECPrivateKey(keyDER)
		if err != nil {
			klog.Errorf("failed to ParseECPrivateKey, error: %v", err)
			return err
		}
		cert, err := x509.ParseCertificate(certDER)
		if err != nil {
			klog.Errorf("failed to parse Certificate %v", err)
			return err
		}
		config.UpdateConfig(nil, nil, certDER, keyDER)
		if err = WriteKeyAndCert(config.Config.TLSPrivateKeyFile, config.Config.TLSCertFile, key, cert); err != nil {
			return fmt.Errorf("failed to save the edge key and certificate to file: %s, error: %v", config.Config, err)
		}
	}
	return nil
}
