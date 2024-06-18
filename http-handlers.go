package main

import (
    "strings"
    "net/http"
    "net/url"
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
    fmt.Fprintf(w, "Page not found")
}

// Examples: /ABCDEF, /1234A6
func codeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    url, err := getUrl(vars["code"])
    if err == sql.ErrNoRows { 
        notFoundHandler(w, r)
        return
    } else if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "Internal error")
        log.Printf("Server error. \n\n Request: \n%v\n Error: \n%v\n", r, err)
        return
    }

    templates.Redirect.Execute(w, url)
}

// Checks the HTTP Requests, parses the urlOriginal, and returns it if found.
// If the request is not correct, error is written inside the function, and "" is returned.
// If the request is correct, the function returns the urlOriginal parameter passed.
func parseCodeRequest(w http.ResponseWriter, r *http.Request) string {
    // Step 1. Check that the Content-Type is application/json (if present)
    ct := r.Header.Get("Content-Type")
    if ct != "" {
        mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
        if mediaType != "application/json" {
            msg := "Content-Type header is not application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return ""
        }
    }
    
    // Step 2. Enforce maximum read of 8KB
    // Too large request body will result in Decode() producing an error.
    r.Body = http.MaxBytesReader(w, r.Body, 8192)
    
    // Step 3. Parse JSON
    var data map[string]interface{}
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return ""
    }
    
    // Step 4. Look for urlOriginal in the json
    urlOriginal, ok := data["urlOriginal"].(string)
    if !ok {
        http.Error(w, "urlOriginal not found in the request body.", http.StatusBadRequest)
        return ""
    }
    
    // Step 5. Fail the request if urlOriginal is bigger than 2048 bytes.
    if len(urlOriginal) > 2048 {
        http.Error(w, "urlOriginal length must be less than 2048 bytes.", http.StatusBadRequest)
        return ""
    }

    // Step 6. Fail the request if urlOriginal is not a URL
    parsedUrl, err := url.ParseRequestURI(urlOriginal)
    if err != nil {
        http.Error(w, "Error in parsing urlOriginal. Please make sure a proper URL is passed.", http.StatusBadRequest)
        return ""
    }
    
    // Step 7. Fail the request if the URL Host is shorturl.space
    if parsedUrl.Hostname() == "shorturl.space" {
        http.Error(w, "urlOriginal can't have shorturl.space as a host", http.StatusBadRequest)
        return ""
    }
    return urlOriginal
}


// /api/new
func newCodeHandler(w http.ResponseWriter, r *http.Request) {
    urlOriginal := parseCodeRequest(w, r);
    if urlOriginal == "" { 
        // the request isn't correct
        return
    }

    code, err := createCode(urlOriginal)
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    
    output := make(map[string]interface{})
    output["urlOriginal"] = urlOriginal
    output["urlCode"] = code

    w.Header().Add("Content-Type", "application/json")
    j, err := json.Marshal(output)
    if err != nil {
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    w.Write(j)
}

// Static handlers for index.html, index.css, index.js

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html");
}

func cssHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.css");
}

func jsHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.js");
}
