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

func NewClient(url string, token string) *Client {
	return &Client{
		url:        url,
		token:      token,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// TODO: routes
// save blob
// find unique blob
// find blobs by bucket
// download blob
// delete blob
// sign url
