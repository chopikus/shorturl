package main

import (
  "os"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "time"
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

func autoRemoveExpired() {
   ticker := time.NewTicker(1 * time.Minute)

   for _ = range ticker.C {
      err := removeExpiredCodes();
      if err != nil {
        log.Printf("[AUTO] Error removing expired codes! %v\n", err)
      } else {
        log.Println("[AUTO] Removed expired codes")
      }
   }
}

func main() {
   go autoRemoveExpired()
   log.Println("autoRemove service started")
   r := NewHandler()
   host := os.Getenv("SHORTURL_SERVER_ADDRESS")
   
   log.Printf("Starting shorturl.space server... SERVER_ADDRESS=%s\n", host)
   log.Fatal(http.ListenAndServe(host, r))
   select {}
}
