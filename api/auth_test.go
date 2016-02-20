package api

import (
	"encoding/base64"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenAuth(t *testing.T) {
	assert := assert.New(t)
	var a AuthProvider = NewTokenAuth("123", "abc")
	assert.True(a.Auth(requestWithAuthHeader("Bearer", "123")))
	assert.True(a.Auth(requestWithAuthHeader("Bearer", "abc")))
	assert.False(a.Auth(requestWithAuthHeader("Bearer", "xyz")))
	assert.False(a.Auth(requestWithAuthHeader("Basic", "abc")))
	assert.False(a.Auth(requestWithAuthHeader("", "")))
}

func TestBasicAuth(t *testing.T) {
	assert := assert.New(t)
	authMap := map[string]string{
		"user": "pass",
		"abc":  "123",
	}
	var a AuthProvider = NewBasicAuth(authMap)
	assert.True(a.Auth(requestWithAuthHeader("Basic", "user:pass")))
	assert.True(a.Auth(requestWithAuthHeader("Basic", "abc:123")))
	assert.False(a.Auth(requestWithAuthHeader("Basic", "abc:abc")))
	assert.False(a.Auth(requestWithAuthHeader("Basic", ":abc")))
	assert.False(a.Auth(requestWithAuthHeader("Basic", "abc:")))
	assert.False(a.Auth(requestWithAuthHeader("Basic", "abc")))
	assert.False(a.Auth(requestWithAuthHeader("Bearer", "user:pass")))
	assert.False(a.Auth(requestWithAuthHeader("Bearer", "")))
}

func TestIPAuth(t *testing.T) {
	assert := assert.New(t)
	var a AuthProvider = NewIPAuth("127.0.0.1", "127.0.1.1")
	assert.True(a.Auth(requestWithRemoteAddr("127.0.0.1")))
	assert.True(a.Auth(requestWithRemoteAddr("127.0.1.1")))
	assert.False(a.Auth(requestWithRemoteAddr("127.0.2.2")))
	assert.False(a.Auth(requestWithRemoteAddr("192.168.0.1")))
	assert.False(a.Auth(requestWithRemoteAddr("")))
}
func requestWithAuthHeader(method string, content string) *http.Request {
	r, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		return nil
	}
	r.Header.Set("Authorization", method+" "+base64.StdEncoding.EncodeToString([]byte(content)))
	return r
}
func requestWithRemoteAddr(ip string) *http.Request {
	r, err := http.NewRequest("GET", "http://127.0.0.1", nil)
	if err != nil {
		return nil
	}
	r.RemoteAddr = ip + ":6000"
	return r
}
