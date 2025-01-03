package handlers

import (
    "fmt"
    "net/http"
)


func ProtectedEmployer(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello there you are in the Employer Protected Zone :)\n")
}
