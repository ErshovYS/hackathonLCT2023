package storage

import (
	"bytes"
	"fmt"
	"invest/internal/config"
	"log"
	"mime/multipart"
	"net/http"
)

type sberStore struct {
	url    string
	origin string
}

func (s *sberStore) UploadFile(filename string, file []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	err := writer.WriteField("key", filename)
	if err != nil {
		return "", err
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}

	_, err = part.Write(file)
	if err != nil {
		return "", err
	}

	contentType := writer.FormDataContentType()
	writer.Close()

	req, err := http.NewRequest("POST", s.url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Origin", s.origin)

	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != http.StatusNoContent {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
		return "", fmt.Errorf("request failed with response code: %d", rsp.StatusCode)
	}

	return fmt.Sprintf("%s/%s", s.url, filename), nil
}

func New(cfg config.StorageConfig) *sberStore {
	return &sberStore{
		url:    cfg.URL,
		origin: cfg.Origin,
	}
}
