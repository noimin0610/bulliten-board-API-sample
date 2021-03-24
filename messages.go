package messages

import (
	// "encoding/json"
	"fmt"
	// "html"
	// "io"
	// "log"
	"net/http"
)

func Messages(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
            fmt.Fprint(w, "Hello World!\n")
    case http.MethodPut:
            http.Error(w, "403 - Forbidden", http.StatusForbidden)
    default:
            http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
    }
}