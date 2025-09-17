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

func (c *Client) setAuth(req *http.Request) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
}
