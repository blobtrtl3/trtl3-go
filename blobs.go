package trtl3

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func (c *Client) UploadBlob(bucket string, file string, r io.Reader) (bool, error) {
	endpoint := fmt.Sprintf("%s/blobs", c.Url)

	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)

	formFile, err := w.CreateFormFile("blob", filepath.Base(file))
	if err != nil {
		return false, fmt.Errorf("error creating form file: %s", err)
	}

	if _, err := io.Copy(formFile, r); err != nil {
		return false, fmt.Errorf("error copying object content: %s", err)
	}

	if err := w.Close(); err != nil {
		return false, fmt.Errorf("error closing writer: %s", err)
	}

	formBucket, err := w.CreateFormField("bucket")
	if err != nil {
		return false, fmt.Errorf("error writing bucket: %s", err)
	}

	formBucket.Write([]byte(bucket))

	req, err := http.NewRequest(http.MethodPost, endpoint, &buffer)
	if err != nil {
		return false, fmt.Errorf("error creating request: %w", err)
	}

  req.Header.Set("Content-Type", w.FormDataContentType())
  req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

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
