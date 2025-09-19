package trtl3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

func (c *Client) UploadBlob(bucket string, r io.Reader) (bool, error) {
	endpoint := fmt.Sprintf("%s/blobs", c.Url)

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)

	formBucket, err := w.CreateFormField("bucket")
	if err != nil {
		return false, fmt.Errorf("error writing bucket: %s", err)
	}

	formBucket.Write([]byte(bucket))

	buf := make([]byte, 512)
	n, err := r.Read(buf)
	if err != nil {
		return false, fmt.Errorf("error reading blob: %s", err)
	}

	var mime = "application/octet-stream"

	if n > 0 {
		mime = http.DetectContentType(buf[:n])
	}

	r = io.MultiReader(bytes.NewReader(buf[:n]), r)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", `form-data; name="blob"; filename="blob"`)
	partHeader.Set("Content-Type", mime)
	formFile, err := w.CreatePart(partHeader)
	if err != nil {
		return false, fmt.Errorf("error creating form file: %s", err)
	}

	if _, err := io.Copy(formFile, r); err != nil {
		return false, fmt.Errorf("error copying blob content: %s", err)
	}

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

	_, err = c.UploadBlob(bucket, r)
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
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Blobs []BlobInfo `json:"blobs"`
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("error deserializing response %w", err)
	}

	return response.Blobs, nil
}

func (c *Client) FindUniqueBlob(bucket, id string) (BlobInfo, error) {
	url := fmt.Sprintf("%s/blobs/%s/%s", c.Url, bucket, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return BlobInfo{}, fmt.Errorf("error trying to create the request: %s", err)
	}

	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return BlobInfo{}, fmt.Errorf("error while doing a request to the server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return BlobInfo{}, fmt.Errorf("failed trying to find created blobs(status: %d)", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return BlobInfo{}, fmt.Errorf("error reading response body: %s", err)
	}

	blobInfo := BlobInfo{}

	if err = json.Unmarshal(bodyBytes, &blobInfo); err != nil {
		return blobInfo, fmt.Errorf("error deserializing response %s", err)
	}

	return blobInfo, nil
}

func (c *Client) SignUrl(sign Signature) (string, error) {
	endpoint := fmt.Sprintf("%s/blobs/sign", c.Url)

	request, err := json.Marshal(sign)
	if err != nil {
		return "", fmt.Errorf("error trying to create the request: %s", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(request))
	if err != nil {
		return "", fmt.Errorf("error trying to create the request: %s", err)
	}

	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error while doing a request to the server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("failed trying to create a signed url(status: %d)", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	var response struct {
		Url string `json:"url"`
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("error deserializing response %s", err)
	}

	return response.Url, nil
}

func (c *Client) DeleteBlob(bucket, id string) (bool, error) {
	endpoint := fmt.Sprintf("%s/blobs/%s/%s", c.Url, bucket, id)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("error trying to create the request: %s", err)
	}

	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error while doing a request to the server: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed trying to delete a blob (status: %d)", res.StatusCode)
	}

	return true, nil
}

func (c *Client) DownloadBlob(bucket, id string) (io.ReadCloser, error) {
	endpoint := fmt.Sprintf("%s/blobs/download/%s/%s", c.Url, bucket, id)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to create the request: %s", err)
	}

	c.setAuth(req)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while doing a request to the server: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("failed trying to download the blob (status: %d)", res.StatusCode)
	}

	return res.Body, nil
}
