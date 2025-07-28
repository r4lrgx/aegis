package main

import (
    "fmt"
    "net/http"
    "os"

    "github.com/r4lrgx/aegis/config"
    "github.com/r4lrgx/aegis/endpoints"
    "github.com/r4lrgx/aegis/middleware"
    "github.com/r4lrgx/aegis/utils"

    "github.com/gorilla/mux"
)

func main() {
    if !utils.ValidateWebhook(config.Webhook) {
        utils.Log("This webhook is invalid or empty. Exiting...")
        os.Exit(1)
    }

    if _, err := os.Stat("IPS.txt"); os.IsNotExist(err) {
        os.WriteFile("IPS.txt", []byte("--------------[github.com/r4lrgx/aegis]--------------"), 0644)
    }

    os.MkdirAll("uploads", os.ModePerm)

    r := mux.NewRouter()
    r.Use(middleware.IPLogger)

    r.HandleFunc("/get", endpoints.GET).Methods("GET")
    r.HandleFunc("/post", endpoints.POST).Methods("POST")
    r.HandleFunc("/delete", endpoints.DELETE).Methods("DELETE")

    addr := fmt.Sprintf(":%d", config.Port)
    utils.Log(fmt.Sprintf("Valid webhook detected. Server is starting on port %d...", config.Port))
    http.ListenAndServe(addr, r)
}
