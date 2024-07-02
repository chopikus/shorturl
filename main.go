package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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
	log.Println("Starting autoRemove service...")
	ticker := time.NewTicker(1 * time.Minute)

	for _ = range ticker.C {
		err := removeExpiredCodes()
		if err != nil {
			log.Printf("[AUTO] Error removing expired codes! %v\n", err)
		} else {
			log.Println("[AUTO] Removed expired codes")
		}
	}
}

func serveHttps() {
	r := NewHandler()
	host := os.Getenv("SHORTURL_HTTPS_ADDRESS")
	certFilePath := os.Getenv("SHORTURL_CERTFILE")
	keyFilePath := os.Getenv("SHORTURL_KEYFILE")
	log.Printf("Starting HTTPS server.. SHORTURL_HTTPS_ADDRESS=%s\n", host)
	log.Fatal(http.ListenAndServeTLS(host, certFilePath, keyFilePath, r))
}

func serveHttp() {
	r := NewHandler()
	host := os.Getenv("SHORTURL_SERVER_ADDRESS")
	log.Printf("Starting HTTP server.. SHORTURL_SERVER_ADDRESS=%s\n", host)
	log.Fatal(http.ListenAndServe(host, r))
}

func main() {
	go autoRemoveExpired()

	//go serveHttps();
	serveHttp()

	select {}
}
