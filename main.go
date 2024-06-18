package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
)

func NewHandler() http.Handler {
    r := mux.NewRouter()
    
    r.HandleFunc("/", indexHandler).
      Methods("GET")

    r.HandleFunc("/index.css", cssHandler).
      Methods("GET")

    r.HandleFunc("/index.js", jsHandler).
      Methods("GET")

    r.HandleFunc("/{code:[1-9A-Z]{6}}", codeHandler).
      Methods("GET")

    r.HandleFunc("/api/new", newCodeHandler).
      Methods("POST")
    
    r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
    r.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)
    
    return r
}

func main() {
   r := NewHandler()
   log.Fatal(http.ListenAndServe("localhost:8000", r))
}
