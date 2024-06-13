package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
)

func codeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Code: %v\n", vars["code"])
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from main page")
}
