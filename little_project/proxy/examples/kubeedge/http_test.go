/*
Copyright 2019 The KubeEdge Authors.

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

package http

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
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"k8s.io/client-go/util/cert"
)

const (
	Path     = "/tmp/kubeedge/testData/"
	BaseName = "edge"
	CertFile = Path + BaseName + ".crt"
	KeyFile  = Path + BaseName + ".key"
	Method   = "GET"
	URL      = "kubeedge.io"
)

// TestNewHttpClient() tests the creation of a new HTTP client
func TestNewHttpClient(t *testing.T) {
	httpClient := NewHTTPClient()
	if httpClient == nil {
		t.Fatal("Failed to build HTTP client")
	}
}

// TestGetCACert
func TestGetCACert(t *testing.T) {
	// kubeedge url
	url := "https://192.168.0.247:10002/ca.crt"

	// edgemesh-server url
	//url := "https://192.168.0.10:10002/ca.crt"
	client := NewHTTPClient()
	req, err := BuildRequest("GET", url, nil, "", "")
	if err != nil {
	}
	res, err := SendRequest(req, client)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	caCert, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(caCert)

	//save the ca.crt to file
	ca, err := x509.ParseCertificate(caCert)
	if err != nil {
		t.Fatalf("failed to parse the CA certificate, error: %v", err)
	}
	t.Log("ca: ", ca)

	block := pem.Block{
		Type: cert.CertificateBlockType,
		Bytes: ca.Raw,
	}
	tmp := pem.EncodeToMemory(&block)
	t.Log(string(tmp))
}

func GetCACert(url string, t *testing.T) []byte {
	client := NewHTTPClient()
	req, err := BuildRequest("GET", url, nil, "", "")
	if err != nil {
	}
	res, err := SendRequest(req, client)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	caCert, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(caCert)
	return caCert
}

func TestGetEdgeCert(t *testing.T) {
	// kubeedge url
	//URL := "https://192.168.0.247:10002/edge.crt"
	//caURL := "https://192.168.0.247:10002/ca.crt"

	// edgemesh server url
	URL := "https://192.168.0.10:10002/agent.crt"
	caURL := "https://192.168.0.10:10002/ca.crt"

	// kubeedge TOKEn
	//token := "78530ed250f06fffd4bf1e95724ac7f503eb1e86cf235096a0e89cada546a9d7.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjMyOTE4MDl9.FiXj97Hw8_tEj0Z4ZiWnBwzsOqy8xkMGkPDolpE57fM"
	// edgemesh TOKEn
	token := "78530ed250f06fffd4bf1e95724ac7f503eb1e86cf235096a0e89cada546a9d7.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjMzMDgwNDl9.hgGcjE-VXwPdYppjTsO0caq8TyEKGhQCHOQPlmXQEfA"
	cacert := GetCACert(caURL, t)

	// validate the CA certificate by hashcode
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 4 {
		t.Fatal("token credentials are in the wrong format")
	}
	ok, hash, newHash := ValidateCACerts(cacert, tokenParts[0])
	if !ok {
		t.Fatalf("failed to validate CA certificate. tokenCAhash: %s, CAhash: %s", hash, newHash)
	}


	caPem := pem.EncodeToMemory(&pem.Block{Bytes: cacert, Type: cert.CertificateBlockType})
	pk, edgeCert, err := GetEdgeCert(URL, caPem, tls.Certificate{}, strings.Join(tokenParts[1:], "."))
	if err != nil {
		t.Fatalf("failed to get edge certificate from the cloudcore, error: %v", err)
	}
	t.Log("pk: ", pk)
	t.Log("edgeCert: ", edgeCert)

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

// GetEdgeCert applies for the certificate from cloudcore
func GetEdgeCert(url string, capem []byte, cert tls.Certificate, token string) (*ecdsa.PrivateKey, []byte, error) {
	pk, csr, err := getCSR()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CSR: %v", err)
	}

	client, err := NewHTTPClientWithCA(capem, cert)
	if err != nil {
		return nil, nil, fmt.Errorf("falied to create http client:%v", err)
	}

	req, err := BuildRequest("GET", url, bytes.NewReader(csr), token, "edgenode")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate http request:%v", err)
	}

	res, err := SendRequest(req, client)
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

func getCSR() (*ecdsa.PrivateKey, []byte, error) {
	CR := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{"kubeEdge"},
			Locality:     []string{"Hangzhou"},
			Province:     []string{"Zhejiang"},
			CommonName:   "kubeedge.io",
		},
	}
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, CR, pk)
	if err != nil {
		return nil, nil, err
	}

	return pk, csr, nil
}

// TestNewHTTPSClient() tests the creation of a new HTTPS client with proper values
func TestNewHTTPSClient(t *testing.T) {
	err := GenerateTestCertificate(Path, BaseName, BaseName)
	if err != nil {
		t.Errorf("Error in generating fake certificates: %w", err)
		return
	}
	certificate, err := tls.LoadX509KeyPair(CertFile, KeyFile)
	if err != nil {
		t.Errorf("Error in loading key pair: %w", err)
		return
	}
	type args struct {
		certFile string
		keyFile  string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Client
		wantErr bool
	}{
		{
			name: "TestNewHTTPSClient: ",
			args: args{
				keyFile:  KeyFile,
				certFile: CertFile,
			},
			want: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:      x509.NewCertPool(),
						Certificates: []tls.Certificate{certificate},
						MinVersion:   tls.VersionTLS12,
						CipherSuites: []uint16{
							tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
						},
						InsecureSkipVerify: true},
				},
				Timeout: connectTimeout,
			},
			wantErr: false,
		},
		{
			name: "Wrong path given while getting HTTPS client",
			args: args{
				keyFile:  "WrongKeyFilePath",
				certFile: "WrongCertFilePath",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHTTPSClient(tt.args.certFile, tt.args.keyFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHTTPSClient() error = %v, expectedError = %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPSClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNewHTTPClientWithCA() tests the creation of a new HTTP using filled capem
func TestNewHTTPClientWithCA(t *testing.T) {
	err := GenerateTestCertificate(Path, BaseName, BaseName)
	if err != nil {
		t.Errorf("Error in generating fake certificates: %w", err)
		return
	}
	capem, err := ioutil.ReadFile(CertFile)
	if err != nil {
		t.Errorf("Error in loading Cert file: %w", err)
		return
	}
	certificate := tls.Certificate{}

	testPool := x509.NewCertPool()
	if ok := testPool.AppendCertsFromPEM(capem); !ok {
		t.Errorf("cannot parse the certificates")
		return
	}

	type args struct {
		capem       []byte
		certificate tls.Certificate
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Client
		wantErr bool
	}{
		{
			name: "TestNewHTTPClientWithCA: ",
			args: args{
				capem:       capem,
				certificate: certificate,
			},
			want: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						RootCAs:            testPool,
						InsecureSkipVerify: false,
						Certificates:       []tls.Certificate{certificate},
					},
				},
				Timeout: connectTimeout,
			},
			wantErr: false,
		},
		{
			name: "Wrong certifcate given when getting HTTP client",
			args: args{
				capem:       []byte{},
				certificate: certificate,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHTTPClientWithCA(tt.args.capem, tt.args.certificate)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHTTPClientWithCA() error = %w, expectedError = %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHTTPClientWithCA() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBuildRequest() tests the process of message building
func TestBuildRequest(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	token := "token"
	nodeName := "name"

	req, err := http.NewRequest(Method, URL, reader)
	if err != nil {
		t.Errorf("Error in creating new http request message: %w", err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("NodeName", nodeName)

	type args struct {
		method   string
		urlStr   string
		body     io.Reader
		token    string
		nodeName string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "TestBuildRequest: ",
			args: args{
				method:   Method,
				urlStr:   URL,
				body:     reader,
				token:    token,
				nodeName: nodeName,
			},
			want:    req,
			wantErr: false,
		},
		{
			name: "NewRequest failure causes BuildRequest failure: ",
			args: args{
				method:   "INVALID\n",
				urlStr:   URL,
				body:     reader,
				token:    token,
				nodeName: nodeName,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildRequest(tt.args.method, tt.args.urlStr, tt.args.body, tt.args.token, tt.args.nodeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuildRequest() error = %w, expectedError = %v", err, tt.wantErr)
				return
			}
			//needed to handle failure testcase because can't deep compare field in nil
			if got == tt.want && err != nil && tt.wantErr == true {
				return
			}
			if !reflect.DeepEqual(got.Header, tt.want.Header) {
				t.Errorf("BuildRequest() Header = %v, want %v", got, tt.want.Header)
			}
			if !reflect.DeepEqual(got.Body, tt.want.Body) {
				t.Errorf("BuildRequest() Body = %v, want %v", got, tt.want.Body)
			}
		})
	}
}

// TestSendRequestFailure() uses fake data and expects function to fail
func TestSendRequestFailure(t *testing.T) {
	httpClient := NewHTTPClient()
	if httpClient == nil {
		t.Fatal("Failed to build HTTP client")
	}

	req, err := http.NewRequest(Method, URL, bytes.NewReader([]byte{}))
	if err != nil {
		t.Errorf("Error in creating new http request message: %w", err)
		return
	}

	resp, respErr := SendRequest(req, httpClient)
	if resp != nil && respErr == nil {
		t.Errorf("Error, response should not come as data is not valid")
	}
}
