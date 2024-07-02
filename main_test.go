package main

import (
	"bytes"
	"encoding/json"
	"github.com/chopikus/shorturl/templates"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewCode(t *testing.T) {
	applicationHandler := NewHandler()
	s := httptest.NewServer(applicationHandler)
	defer s.Close()

	type TestCase struct {
		description    string
		contentType    string
		body           string
		wantStatusCode int
	}

	testCases := []TestCase{
		{
			description:    "Wrong Content-Type (text/html)",
			contentType:    "text/html",
			body:           `{"urlOriginal": "https://example.com"}`,
			wantStatusCode: 415,
		},
		{
			description:    "Correct Content-Type (empty)",
			contentType:    "",
			body:           `{"urlOriginal": "https://example.com"}`,
			wantStatusCode: 200,
		},
		{
			description:    "Empty body",
			contentType:    "application/json",
			body:           "",
			wantStatusCode: 400,
		},
		{
			description:    "Too large body (20k characters)",
			contentType:    "application/json",
			body:           strings.Repeat("a", 20000),
			wantStatusCode: 400,
		},
		{
			description:    "Wrong JSON key",
			contentType:    "application/json",
			body:           `{"hello": "world"}`,
			wantStatusCode: 400,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			bodyReader := strings.NewReader(testCase.body)

			res, err := http.Post(s.URL+"/api/new", testCase.contentType, bodyReader)

			if err != nil {
				t.Fatalf("%s", err)
			}

			assert.Equal(t, testCase.wantStatusCode, res.StatusCode)
		})
	}
}

// Creates new code using /api/new, returns the code string
func createNewCode(t *testing.T, s *httptest.Server, url string) string {
	bodyReader := strings.NewReader(`{"urlOriginal": "` + url + `"}`)

	res, err := http.Post(s.URL+"/api/new", "application/json", bodyReader)
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
