package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/r4lrgx/aegis/config"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

var cyan = color.New(color.FgCyan)

func Log(message string) {
	now := time.Now()
	timestamp := fmt.Sprintf("[%02d:%02d:%02d]", now.Hour(), now.Minute(), now.Second())
	cyan.Printf("%s %s\n", timestamp, message)
}

func ValidateWebhook(url string) bool {
	if url == "" {
		return false
	}

	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false
	}

	if _, ok := data["id"]; !ok {
		return false
	}

	return true
}

func SendMultipart(writer *multipart.Writer, body *bytes.Buffer, files []string) ([]byte, error) {
	req, err := http.NewRequest("POST", config.Webhook, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reach the webhook")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read webhook response")
	}

	for _, f := range files {
		os.Remove(f)
	}

	return respBody, nil
}

func AttachFile(writer *multipart.Writer, fieldName, path, displayName string) error {
	part, err := writer.CreateFormFile(fieldName, displayName)
	if err != nil {
		return fmt.Errorf("error creating form file")
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error reopening file")
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error writing file to multipart")
	}
	return nil
}

func GetFieldName(i int) string {
	if i == 0 {
		return "file"
	}
	return fmt.Sprintf("file%d", i)
}