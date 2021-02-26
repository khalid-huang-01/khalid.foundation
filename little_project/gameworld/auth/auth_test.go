package auth

import (
	"fmt"
	"testing"
	"time"
)

func TestDefaultApiAuthenticator_AuthForApiRequest(t *testing.T) {
	baseURL := "/users/123/address"
	appID := "1000"
	password := "1000"
	timestamp := time.Now().Unix()
	token := generateToken(baseURL, appID, password, timestamp)
	// 客户端拼接请求并发送
	request := &ApiRequest{
		baseURL:   baseURL,
		token:     token,
		appID:     appID,
		timestamp: timestamp,
	}
	// 服务端接收到请求，进行验证
	authenticator := NewDefaultApiAuthenticator()
	err := authenticator.AuthForApiRequest(request)
	if err != nil {
		t.Log("err: ", err)
	} else {
		t.Log("pass")
	}
}

func TestDefaultApiAuthenticator_AuthForUrl(t *testing.T) {
	// 客户端请求
	baseURL := "/users/123/address"
	appID := "1000"
	password := "1000"
	timestamp := time.Now().Unix()
	token := generateToken(baseURL, appID, password, timestamp)
	url := fmt.Sprintf("%s?appid=%s&timestamp=%d&token=%s", baseURL, appID, timestamp, token)
	// 服务端接收
	request := newApiRequestFromFullURL(url)
	// 进行验证
	authenticator := NewDefaultApiAuthenticator()
	err := authenticator.AuthForApiRequest(request)
	if err != nil {
		t.Log("err: ", err)
	} else {
		t.Log("pass")
	}
}