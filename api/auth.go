package api

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type AuthProvider interface {
	Auth(r *http.Request) bool
}

func NewBasicAuth(auth map[string]string) BasicAuth {
	if auth != nil {
		return BasicAuth(auth)
	}
	return make(BasicAuth)
}

type BasicAuth map[string]string

func (a BasicAuth) Auth(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	method, content, valid := decodeAuthHeader(authHeader)
	if !valid {
		return false
	}
	if method != "Basic" {
		return false
	}
	userpass := strings.SplitN(content, ":", 2)
	if len(userpass) != 2 {
		return false
	}
	if pass, ok := a[userpass[0]]; !ok || pass != userpass[1] {
		return false
	}
	return true
}

func NewTokenAuth(token ...string) *TokenAuth {
	a := &TokenAuth{
		tokens: make(map[string]struct{}),
	}
	for _, t := range token {
		a.tokens[t] = struct{}{}
	}
	return a
}

type TokenAuth struct {
	tokens map[string]struct{}
}

func (a *TokenAuth) Auth(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	method, content, valid := decodeAuthHeader(authHeader)
	if !valid {
		return false
	}
	if method != "Bearer" {
		return false
	}

	if _, ok := a.tokens[content]; ok {
		return true
	}
	return false
}

func decodeAuthHeader(header string) (method string, content string, valid bool) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	dec, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", false
	}
	return parts[0], string(dec), true
}

func authHelper(token string) string {
	return "Bearer " + base64.StdEncoding.EncodeToString([]byte(token))
}

type IPAuth map[string]struct{}

func NewIPAuth(ips ...string) IPAuth {
	a := make(IPAuth)
	for _, ip := range ips {
		a[strings.TrimSpace(ip)] = struct{}{}
	}
	return a
}

func (a IPAuth) Auth(r *http.Request) bool {
	parts := strings.SplitN(r.RemoteAddr, ":", 2)
	if len(parts) == 0 {
		return false
	}

	if _, ok := a[parts[0]]; ok {
		return true
	}
	return false
}
