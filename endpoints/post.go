package endpoints

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/r4lrgx/aegis/utils"
)

func POST(w http.ResponseWriter, r *http.Request) {
	utils.Log("POST request received")

	contentType := r.Header.Get("Content-Type")
	isJSON := strings.HasPrefix(contentType, "application/json")

	var body []byte
	var err error

	if isJSON {
		body, err = HandleJSONPayload(r)
	} else {
		body, err = HandleMultipartPayload(r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func HandleJSONPayload(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON body")
	}

	var payload map[string]interface{}
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, fmt.Errorf("invalid JSON body")
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	var uploadedFiles []string

	if attachments, ok := payload["attachments"].([]interface{}); ok {
		for i, item := range attachments {
			a, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			fileName, _ := a["filename"].(string)
			dataStr, _ := a["data"].(string)

			data, err := base64.StdEncoding.DecodeString(dataStr)
			if err != nil {
				continue
			}

			tempPath := filepath.Join("uploads", fmt.Sprintf("json_%d_%s", i, fileName))
			if err := os.WriteFile(tempPath, data, 0644); err != nil {
				continue
			}
			uploadedFiles = append(uploadedFiles, tempPath)

			fieldName := utils.GetFieldName(i)
			if err := utils.AttachFile(writer, fieldName, tempPath, fileName); err != nil {
				return nil, err
			}
		}
	}

	delete(payload, "attachments")
	payloadJSON, _ := json.Marshal(payload)
	writer.WriteField("payload_json", string(payloadJSON))
	writer.Close()

	return utils.SendMultipart(writer, &b, uploadedFiles)
}

func HandleMultipartPayload(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(64 << 20)
	if err != nil {
		return nil, fmt.Errorf("failed to parse multipart form")
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	var uploadedFiles []string

	var files []*multipart.FileHeader
	for i := 0; i < 100; i++ {
		key := utils.GetFieldName(i)
		if f := r.MultipartForm.File[key]; len(f) > 0 {
			files = append(files, f...)
		}
	}

	for i, fh := range files {
		fieldName := utils.GetFieldName(i)

		src, err := fh.Open()
		if err != nil {
			return nil, fmt.Errorf("unable to read uploaded file")
		}
		defer src.Close()

		tempPath := filepath.Join("uploads", fmt.Sprintf("form_%d_%s", i, fh.Filename))
		dst, err := os.Create(tempPath)
		if err != nil {
			return nil, fmt.Errorf("error saving file")
		}
		_, err = io.Copy(dst, src)
		dst.Close()
		if err != nil {
			return nil, fmt.Errorf("error copying file")
		}
		uploadedFiles = append(uploadedFiles, tempPath)

		if err := utils.AttachFile(writer, fieldName, tempPath, fh.Filename); err != nil {
			return nil, err
		}
	}

	if payload := r.FormValue("payload_json"); payload != "" {
		writer.WriteField("payload_json", payload)
	}
	writer.Close()

	return utils.SendMultipart(writer, &b, uploadedFiles)
}
