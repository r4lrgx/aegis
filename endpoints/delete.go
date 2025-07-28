package endpoints

import (
    "net/http"
    "github.com/r4lrgx/aegis/utils"
)

func DELETE(w http.ResponseWriter, r *http.Request) {
    utils.Log("DELETE request, for deletion")

    w.Write([]byte("Nigga, you can't"))
}
