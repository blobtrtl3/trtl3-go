package trtl3

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Url        string
	Token      string
	httpClient *http.Client
}

type BlobInfo struct {
	ID        string    `json:"id"`
	Bucket    string    `json:"bucket"`
	Mime      string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type Signature struct {
	Bucket string
	ID     string
	TTL    int
	Once   bool
}

func NewDefaultClient() *Client {
	return &Client{
		Url:        "http://localhost:7713",
		Token:      "trtl3",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func NewClient(url string, token string, timeout time.Duration) *Client {
	return &Client{
		Url:        url,
		Token:      token,
		httpClient: &http.Client{Timeout: timeout * time.Second},
	}
}

func (c *Client) setAuth(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
}
