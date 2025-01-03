package handlers

import (
	"fmt"
	"net/http"
	//"github.com/gorilla/mux"
)

func LandingPage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello there")
}
