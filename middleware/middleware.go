package middleware

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/r4lrgx/aegis/config"
	"github.com/r4lrgx/aegis/utils"
)

type visitor struct {
	lastSeen time.Time
	tokens   int
}

var (
	mu       sync.Mutex
	visitors = make(map[string]*visitor)
)

func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)

		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		v, exists := visitors[ip]

		if !exists || now.Sub(v.lastSeen) > config.RateLimitWindow {
			v = &visitor{
				lastSeen: now,
				tokens:   config.RateLimitMax,
			}
			visitors[ip] = v
		}

		if v.tokens <= 0 {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			utils.Log("Rate limit hit for IP: " + ip)
			return
		}

		v.tokens--
		v.lastSeen = now

		next.ServeHTTP(w, r)
	})
}

func IPLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)

		now := time.Now().Format("[07-30-2025 15:05:05]")
		entry := fmt.Sprintf("%s %s\n", now, ip)

		f, err := os.OpenFile("ips.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString(entry)
		}

		next.ServeHTTP(w, r)
	})
}

func GetIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	ip = strings.TrimPrefix(ip, "[::ffff:")
	ip = strings.TrimSuffix(ip, "]")
	return ip
}
