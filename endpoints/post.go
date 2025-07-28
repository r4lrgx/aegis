package endpoints

import (
    "bytes"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"

    "github.com/r4lrgx/aegis/utils"
    "github.com/r4lrgx/aegis/config"
)

func POST(w http.ResponseWriter, r *http.Request) {
    utils.Log("POST request, files uploaded")

    var b bytes.Buffer
    writer := multipart.NewWriter(&b)

    r.ParseMultipartForm(32 << 20)
    files := r.MultipartForm.File["file"]

    for i, fileHeader := range files {
        file, err := fileHeader.Open()
        if err != nil {
            http.Error(w, "Unable to read uploaded file", http.StatusInternalServerError)
            return
        }
        defer file.Close()

        dst, err := os.Create(filepath.Join("uploads", fileHeader.Filename))
        if err != nil {
            http.Error(w, "Error saving file", http.StatusInternalServerError)
            return
        }
        io.Copy(dst, file)
        dst.Close()

        fw, _ := writer.CreateFormFile("file"+string(rune(i+1)), fileHeader.Filename)
        fr, _ := os.Open(filepath.Join("uploads", fileHeader.Filename))
        io.Copy(fw, fr)
        fr.Close()
    }

    if r.FormValue("payload_json") != "" {
        writer.WriteField("payload_json", r.FormValue("payload_json"))
    }
    writer.Close()

    req, _ := http.NewRequest("POST", config.Webhook, &b)
    req.Header.Set("Content-Type", writer.FormDataContentType())

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        http.Error(w, "Failed to reach the webhook", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    for _, fileHeader := range files {
        os.Remove(filepath.Join("uploads", fileHeader.Filename))
    }

    w.Write(body)
}
