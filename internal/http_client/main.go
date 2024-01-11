package http_client

import (
	"fmt"
	"net/http"
)

type AuthTransport struct {
    Token string
    Transport http.RoundTripper
}

var defaultClient *http.Client
var defaultToken string

func New(token string) *http.Client {
	t := &AuthTransport{
		Token: token,
		Transport: http.DefaultTransport,
	}
	if defaultClient == nil {
		defaultClient = &http.Client{Transport: t}
	}
	
	return defaultClient
}

func Default() *http.Client {
	return defaultClient
}

func DefaultToken() string {
	return defaultToken
}

func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	defaultToken = t.Token
    req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.Token))
    return t.Transport.RoundTrip(req)
}
