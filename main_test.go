package main

import (
    "bytes"
    "encoding/json"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "github.com/chopikus/shorturl/templates"
)


// Creates new code using /api/new, returns the code string
func createNewCode(t *testing.T, s *httptest.Server, url string) string {
    bodyReader := strings.NewReader(`{"urlOriginal": "` + url + `"}`)
   
    res, err := http.Post(s.URL + "/api/new", "application/json", bodyReader)
   
    if err != nil {
        t.Fatalf("%s", err)
    }

    defer res.Body.Close()
    response, err := io.ReadAll(res.Body)

    if err != nil {
        t.Fatalf("%s", err)
    }
 
    j := make(map[string]any)
    if err = json.Unmarshal(response, &j); err != nil {
        t.Logf("%s\n", response)
        t.Fatalf("%s", err)
    }

    return j["urlCode"].(string)
}

func ioReaderCheck(t *testing.T, r io.Reader, expected bytes.Buffer) bool {
    rBytes, err := io.ReadAll(r)

    if err != nil {
        t.Fatal(err)
        return false
    }
     
    return bytes.Equal(rBytes, expected.Bytes())
}

func TestCodeHandler(t *testing.T) {
    const url = "https://example.com"

    // Step 1. Create the server
    applicationHandler := NewHandler()
    s := httptest.NewServer(applicationHandler)
    defer s.Close()

    // Step 2. Generate new code for url
    code := createNewCode(t, s, url)
    
    // Step 3. Get the response for the code
    res, err := http.Get(s.URL + "/" + code)
    if err != nil {
        t.Fatalf("%s", err)
    }
    defer res.Body.Close()

    var expected bytes.Buffer
    templates.Redirect.Execute(&expected, url)

    if !ioReaderCheck(t, res.Body, expected) {
        t.Fatalf("!ioReaderCheck, res.Body=%v \n\n\n expected=%v", res.Body, expected)
    }
}
