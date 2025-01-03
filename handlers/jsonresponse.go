package handlers

import (
     "encoding/json"
     "net/http"
)

func JsonResponse(w http.ResponseWriter, contentType string, content string, status int, response map[string]interface{}) {
    w.Header().Set(contentType, content)
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}
