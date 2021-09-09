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
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"khalid.jobs/caserver/config"
	"net"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"

)

// SignCerts creates server's certificate and key
func SignCerts() ([]byte, []byte, error) {
	cfg := &certutil.Config{
		CommonName:   "KubeEdge",
		Organization: []string{"KubeEdge"},
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		AltNames: certutil.AltNames{
			DNSNames: config.Config.DNSNames,
			IPs:      getIps(config.Config.AdvertiseAddress),
		},
	}

	certDER, keyDER, err := NewCloudCoreCertDERandKey(cfg)
	if err != nil {
		return nil, nil, err
	}

	return certDER, keyDER, nil
}

func getIps(advertiseAddress []string) (Ips []net.IP) {
	for _, addr := range advertiseAddress {
		Ips = append(Ips, net.ParseIP(addr))
	}
	return
}

// GenerateToken will create a token consisting of caHash and jwt Token and save it to secret
func GenerateToken() error {
	expiresAt := time.Now().Add(time.Hour * 24).Unix()

	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = jwt.StandardClaims{
		ExpiresAt: expiresAt,
	}

	keyPEM := getCaKey()
	tokenString, err := token.SignedString(keyPEM)

	if err != nil {
		return fmt.Errorf("Failed to generate the token for EdgeCore register, err: %v", err)
	}

	caHash := getCaHash()
	// combine caHash and tokenString into caHashAndToken
	caHashToken := strings.Join([]string{caHash, tokenString}, ".")

	klog.Infof("server token: %s", caHashToken)

	klog.Info("Succeed to creating token")
	return nil
}


// getCaHash gets ca-hash
func getCaHash() string {
	caDER := config.Config.Ca
	digest := sha256.Sum256(caDER)
	return hex.EncodeToString(digest[:])
}

// getCaKey gets caKey to encrypt token
func getCaKey() []byte {
	caKey := config.Config.CaKey
	return caKey
}
