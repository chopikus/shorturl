package main

import (
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func createNewCode(t *testing.T, s *httptest.Server) {
    bodyReader := strings.NewReader(`{"urlOriginal": "https://chopikus.github.io"}`)
   
    res, err := http.Post(s.URL + "/api/new", "application/json", bodyReader)
   
    if err != nil {
        t.Fatalf("%s", err)
    }

    defer res.Body.Close()
    response, err := io.ReadAll(res.Body)

    if err != nil {
        t.Fatalf("%s", err)
    }

    t.Logf("%s", response)
}

func TestCodeHandler(t *testing.T) {
    /* Step 1. Create new code */
       
    applicationHandler := NewHandler()
    
    s := httptest.NewServer(applicationHandler)
    defer s.Close()

    createNewCode(t, s)
}
