package endpoints

import (
    "encoding/json"
    "net/http"

    "github.com/r4lrgx/aegis/config"
    "github.com/r4lrgx/aegis/utils"
)

func GET(w http.ResponseWriter, r *http.Request) {
    utils.Log("GET request, keys replaced")

    resp, err := http.Get(config.Webhook)
    if err != nil {
        http.Error(w, "Unable to contact the webhook", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&data)

    for _, key := range []string{"id", "token", "channel_id", "guild_id", "url"} {
        data[key] = "Aegis_" + utils.RandomString(8)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
