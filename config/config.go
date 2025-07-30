package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	Webhook         string
	Port            int
	RateLimitWindow time.Duration
	RateLimitMax    int
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	Webhook = os.Getenv("WEBHOOK")
	if Webhook == "" {
		log.Fatal("WEBHOOK is required in .env")
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port <= 0 {
		log.Fatal("Invalid or missing PORT")
	}
	Port = port

	mins, err := strconv.Atoi(os.Getenv("RATE_LIMIT_WINDOW_MINUTES"))
	if err != nil || mins <= 0 {
		log.Fatal("Invalid or missing RATE_LIMIT_WINDOW_MINUTES")
	}
	RateLimitWindow = time.Duration(mins) * time.Minute

	max, err := strconv.Atoi(os.Getenv("RATE_LIMIT_MAX"))
	if err != nil || max <= 0 {
		log.Fatal("Invalid or missing RATE_LIMIT_MAX")
	}
	RateLimitMax = max
}
