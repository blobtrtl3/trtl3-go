package trtl3

import (
	"net/http"
	"time"
)

type Client struct {
	url        string
	token      string
	httpClient *http.Client
}

func New(url string, token string) *Client {
	return &Client{
		url:        url,
		token:      token,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}
