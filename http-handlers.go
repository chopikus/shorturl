package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "log"
    "database/sql"
    "context"
    "github.com/chopikus/url-shortener/html"
)

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "Page not found. We can't seem to find the page you're looking for.")
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    url, err := getUrl(vars["code"])
    if err == sql.ErrNoRows { 
        notFoundHandler(w, r)
        return
    } else if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Internal error")
        log.Println("Server error. Request: %v\n, Error: %v\n", err)
        return
    }

    component := html.Redirect(url)
    component.Render(context.Background(), w)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
    component := html.MainPage()
	component.Render(context.Background(), w)
}
