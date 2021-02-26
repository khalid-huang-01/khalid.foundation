package auth

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultExpiredTimeInterval int64 = 2
)

type AuthToken struct {
	token string
	createTime int64
	expiredTimeInterval int64 //秒
}

func (a *AuthToken) getToken() string {
	return a.token
}

func (a *AuthToken) isExpired() bool {
	now := time.Now()
	createTime := time.Unix(a.createTime, 0)
	expireTime := createTime.Add(time.Duration(a.expiredTimeInterval) * time.Second)
	return now.After(expireTime)
}

func (a *AuthToken) match(token string) bool {
	return a.token == token
}

// 这个方法理论上只有客户端会调用，服务端应该只是直接注入token就可以了(在客户端生成token)
func generateToken(url string, appID string, passwd string, createTime int64) string {
	str := fmt.Sprintf("%s@%s@%s@%d", url,
		appID, passwd, createTime)
	fmt.Println(str)
	tmp := md5.Sum([]byte(str))
	token := fmt.Sprintf("%x", tmp)
	return token
}


// 代表一个用户请求
type ApiRequest struct {
	baseURL string
	token string
	appID string
	timestamp int64
}

// 简单处理，接收一个完整的请求，从中解析出信息
// /users/khalid/address?appid=sss&timestamp=xxxx&token=xxxx
func newApiRequestFromFullURL(url string) *ApiRequest {
	request := &ApiRequest{}
	parts := strings.Split(url, "?")
	request.baseURL = parts[0]
	params := make(map[string]string, 0)
	queries := strings.Split(parts[1], "&")
	for _, part := range queries {
		kv := strings.Split(part, "=")
		params[kv[0]] = kv[1]
	}
	request.appID = params["appid"]
	request.token = params["token"]
	tmp, _ := strconv.Atoi(params["timestamp"])
	request.timestamp = int64(tmp)
	return request
}

type CredentialStorage interface {
	GetPasswordByAppID(string) string
}

type MapCredentialStorage struct {
	db map[string]string
}

func (s *MapCredentialStorage) GetPasswordByAppID(appID string) string {
	return s.db[appID]
}


type ApiAuthenticator interface {
	AuthForUrl(url string)
	AuthForApiRequest(request ApiRequest)
}

type DefaultApiAuthenticator struct {
	credentialStorage CredentialStorage
}

func NewDefaultApiAuthenticator() *DefaultApiAuthenticator {
	return &DefaultApiAuthenticator{
		credentialStorage: &MapCredentialStorage{db: map[string]string{
			"1000": "1000",
			"1001": "1001",
			"1002": "1002",
		}},
	}
}

// 模拟接收到url之后怎么做
func (a *DefaultApiAuthenticator) AuthForUrl(url string) error {
	request := newApiRequestFromFullURL(url)
	return a.AuthForApiRequest(request)
}

func (a *DefaultApiAuthenticator) AuthForApiRequest(request *ApiRequest) error {
	// 判断是否过期
	token := AuthToken{
		token:               request.token,
		createTime:          request.timestamp,
		expiredTimeInterval: DefaultExpiredTimeInterval,
	}
	if token.isExpired() {
		return fmt.Errorf("token is expired")
	}

	// 获取密码，生产token
	password := a.credentialStorage.GetPasswordByAppID(request.appID)
	serverAuthTokenStr := generateToken(request.baseURL, request.appID, password, request.timestamp)

	// 判断token是否一致
	if !token.match(serverAuthTokenStr) {
		return fmt.Errorf("token is not match")
	}
	return nil

}


