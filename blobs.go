package trtl3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (c *Client) UploadBlob(bucket string, name string, r io.Reader) (bool, error) {
	endpoint := fmt.Sprintf("%s/blobs", c.Url)

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)

	formFile, err := w.CreateFormFile("blob", name)
	if err != nil {
		return false, fmt.Errorf("error creating form file: %s", err)
	}

	if _, err := io.Copy(formFile, r); err != nil {
		return false, fmt.Errorf("error copying blob content: %s", err)
	}

	formBucket, err := w.CreateFormField("bucket")
	if err != nil {
		return false, fmt.Errorf("error writing bucket: %s", err)
	}

	formBucket.Write([]byte(bucket))

	if err := w.Close(); err != nil {
		return false, fmt.Errorf("error closing writer: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, &buffer)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error while doing request to server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return false, fmt.Errorf("upload failed with status (%d)", res.StatusCode)
	}

	return true, nil
}

func (c *Client) UploadBlobByPath(bucket string, path string) (bool, error) {
	r, err := os.Open(path)
	if err != nil {
		return false, err
	}

	_, err = c.UploadBlob(bucket, filepath.Base(path), r)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) FindBlobsByBucket(bucket string) ([]BlobInfo, error) {
	url := fmt.Sprintf("%s/blobs?bucket=%s", c.Url, bucket)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to create the request: %s", err)
	}

	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing a request to the server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed trying to find created blobs(status: %d)", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var response struct {
		Blobs []BlobInfo `json:"blobs"`
	}

	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, fmt.Errorf("error deserializing response %w", err)
	}

	return response.Blobs, nil
}
