package main

import (
    "net/http"
    "github.com/gorilla/mux"
    "fmt"
    "log"
    "database/sql"
    "github.com/chopikus/shorturl/templates"
    "encoding/json"
)

func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusMethodNotAllowed)
    fmt.Fprintf(w, "Method not allowed!")
}

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

    templates.Redirect.Execute(w, url)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
    templates.Index.Execute(w, nil)
}

func generateCodeHandler(w http.ResponseWriter, r *http.Request) {
    // https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
    var data map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    originalUrl, ok := data["urlOriginal"].(string)
    if !ok {
        http.Error(w, "please put urlOriginal in the body", http.StatusBadRequest)
        return
    }

    code, err := createCode(originalUrl)
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    
    output := make(map[string]interface{})
    output["urlOriginal"] = originalUrl
    output["urlCode"] = code

    w.Header().Add("Content-Type", "application/json")
    j, err := json.Marshal(output)
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    w.Write(j)
}
