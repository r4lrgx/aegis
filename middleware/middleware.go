package middleware

import (
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"
)

func IPLogger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.Header.Get("X-Forwarded-For")
        if ip == "" {
            ip = r.RemoteAddr
        }

        if strings.HasPrefix(ip, "[::ffff:") {
            ip = strings.TrimPrefix(ip, "[::ffff:")
        }

        now := time.Now()
        entry := fmt.Sprintf(
            "\n[%02d:%02d:%04d:%02d:%02d:%02d]: %s",
            now.Day(), now.Month(), now.Year(),
            now.Hour(), now.Minute(), now.Second(),
            ip,
        )

        f, _ := os.OpenFile("IPS.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        defer f.Close()
        f.WriteString(entry)

        next.ServeHTTP(w, r)
    })
}
