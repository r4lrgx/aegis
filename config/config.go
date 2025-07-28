package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var Webhook string
var Port = 3000
var RateLimitWindow = 15 * time.Minute
var RateLimitMax = 100

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not loaded")
	}

	Webhook = os.Getenv("WEBHOOK")
	if Webhook == "" {
		log.Fatal("WEBHOOK not set in .env aborting startup")
	}
}
