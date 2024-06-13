package main

import (
  "database/sql"
  _ "github.com/lib/pq"
  "log"
  "math/rand"
  "time"
  "net/http"
  "github.com/gorilla/mux"
)

var db *sql.DB

func init() {
    /* rand is used to generate new URL codes */
    rand.Seed(time.Now().UnixNano())

    tmpDB, err := sql.Open("postgres", "dbname=urldb user=ihor password=ihor host=/tmp sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    db = tmpDB
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", mainPageHandler)
    r.HandleFunc("/{code:[1-9A-Z]{6}}", codeHandler)
    r.HandleFunc("/request-code", mainPageHandler).
      Methods("POST")
    r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
    r.MethodNotAllowedHandler = http.HandlerFunc(notFoundHandler)
    log.Fatal(http.ListenAndServe("localhost:8000", r))
}
