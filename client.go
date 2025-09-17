package trtl3

import (
	"net/http"
	"time"
)

type Client struct {
	Url        string
	Token      string
	httpClient *http.Client
}

func NewClient(url string, token string) *Client {
	return &Client{
		Url:        url,
		Token:      token,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// TODO: routes
// find unique blob
// find blobs by bucket
// download blob
// delete blob
// sign url
