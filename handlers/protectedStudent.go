package handlers


import (
    "fmt"
    "net/http"
)

func ProtectedStudent(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "You are in the protected Student Zone :)\n")
}
