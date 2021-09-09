package acl

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/edge/pkg/edgehub/common/certutil"
	"github.com/kubeedge/kubeedge/edge/pkg/edgehub/common/http"
	"k8s.io/client-go/util/cert"
	"k8s.io/klog/v2"
)


type ACLManager struct {
	CR *x509.CertificateRequest

	caFile   string
	certFile string
	keyFile  string

	token string
	// Set to time.Now but can be stubbed out for testing
	now func() time.Time

	caURL   string
	certURL string
}

func getIPs(advertiseAddress []string) []net.IP {
	IPs := make([]net.IP, len(advertiseAddress))
	for i, v := range advertiseAddress {
		IPs[i] = net.ParseIP(v)
	}
	return IPs
}

// NewACLManager creates a ACLManager for edge acl management according to EdgeHub config
func NewACLManager(tunnel TunnelACLConfig) *ACLManager {
	certReq := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{"kubeEdge"},
			Locality:     []string{"Hangzhou"},
			Province:     []string{"Zhejiang"},
			CommonName:   "kubeedge.io",
		},
		DNSNames:    []string{"localhost"},
		IPAddresses: getIPs([]string{"127.0.0.1"}),
	}
	return &ACLManager{
		token:    tunnel.Token,
		CR:       certReq,
		caFile:   tunnel.TLSCAFile,
		certFile: tunnel.TLSCertFile,
		keyFile:  tunnel.TLSPrivateKeyFile,
		now:      time.Now,
		caURL:    tunnel.HTTPServer + constants.DefaultCAURL,
		certURL:  tunnel.HTTPServer + constants.DefaultCertURL,
	}
}

// Start starts the ACLManager
func (cm *ACLManager) Start() {
	_, err := cm.getCurrent()
	if err != nil {
		err = cm.applyCerts()
		if err != nil {
			klog.Fatalf("Error: %v", err)
		}
	}
}

// getCurrent returns current edge certificate
func (cm *ACLManager) getCurrent() (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair(cm.certFile, cm.keyFile)
	if err != nil {
		return nil, err
	}
	certs, err := x509.ParseCertificates(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("unable to parse certificate data: %v", err)
	}
	cert.Leaf = certs[0]
	return &cert, nil
}

// applyCerts realizes the certificate application by token
func (cm *ACLManager) applyCerts() error {
	cacert, err := GetCACert(cm.caURL)
	if err != nil {
		return fmt.Errorf("failed to get CA certificate, err: %v", err)
	}

	// validate the CA certificate by hashcode
	tokenParts := strings.Split(cm.token, ".")
	if len(tokenParts) != 4 {
		return fmt.Errorf("token credentials are in the wrong format")
	}
	ok, hash, newHash := ValidateCACerts(cacert, tokenParts[0])
	if !ok {
		return fmt.Errorf("failed to validate CA certificate. tokenCAhash: %s, CAhash: %s", hash, newHash)
	}

	// save the ca.crt to file
	ca, err := x509.ParseCertificate(cacert)
	if err != nil {
		return fmt.Errorf("failed to parse the CA certificate, error: %v", err)
	}

	if err = certutil.WriteCert(cm.caFile, ca); err != nil {
		return fmt.Errorf("failed to save the CA certificate to file: %s, error: %v", cm.caFile, err)
	}

	// get the edge.crt
	caPem := pem.EncodeToMemory(&pem.Block{Bytes: cacert, Type: cert.CertificateBlockType})
	pk, edgeCert, err := cm.GetEdgeCert(cm.certURL, caPem, tls.Certificate{}, strings.Join(tokenParts[1:], "."))
	if err != nil {
		return fmt.Errorf("failed to get edge certificate from the cloudcore, error: %v", err)
	}

	// save the edge.crt to the file
	crt, _ := x509.ParseCertificate(edgeCert)
	if err = certutil.WriteKeyAndCert(cm.keyFile, cm.certFile, pk, crt); err != nil {
		return fmt.Errorf("failed to save the edge key and certificate to file: %s, error: %v", cm.certFile, err)
	}

	return nil
}

// applyCerts realizes the acl application by token
func (cm *ACLManager) generateKey() error {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key, error: %v", err)
	}
	// save the private key to the file
	if err = certutil.WriteKey(cm.keyFile, pk); err != nil {
		return fmt.Errorf("failed to save the private key %s, error: %v", cm.keyFile, err)
	}
	return nil
}

// getCA returns the CA in pem format.
func (cm *ACLManager) getCA() ([]byte, error) {
	return ioutil.ReadFile(cm.caFile)
}


// GetEdgeCert applies for the certificate from cloudcore
func (cm *ACLManager) GetEdgeCert(url string, capem []byte, cert tls.Certificate, token string) (*ecdsa.PrivateKey, []byte, error) {
	pk, csr, err := cm.getCSR()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CSR: %v", err)
	}

	client, err := http.NewHTTPClientWithCA(capem, cert)
	if err != nil {
		return nil, nil, fmt.Errorf("falied to create http client:%v", err)
	}

	nodeName := "edge-node"

	req, err := http.BuildRequest("GET", url, bytes.NewReader(csr), token, nodeName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate http request:%v", err)
	}

	res, err := http.SendRequest(req, client)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != 200 {
		return nil, nil, fmt.Errorf(string(content))
	}

	return pk, content, nil
}

func (cm *ACLManager) getCSR() (*ecdsa.PrivateKey, []byte, error) {
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, cm.CR, pk)
	if err != nil {
		return nil, nil, err
	}

	return pk, csr, nil
}

// ValidateCACerts validates the CA certificate by hash code
func ValidateCACerts(cacerts []byte, hash string) (bool, string, string) {
	if len(cacerts) == 0 && hash == "" {
		return true, "", ""
	}

	newHash := hashCA(cacerts)
	return hash == newHash, hash, newHash
}

func hashCA(cacerts []byte) string {
	digest := sha256.Sum256(cacerts)
	return hex.EncodeToString(digest[:])
}

// GetCACert gets the cloudcore CA certificate
func GetCACert(url string) ([]byte, error) {
	client := http.NewHTTPClient()
	req, err := http.BuildRequest("GET", url, nil, "", "")
	if err != nil {
		return nil, err
	}
	res, err := http.SendRequest(req, client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	caCert, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return caCert, nil
}
