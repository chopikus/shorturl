package main

import (
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCodeHandler(t *testing.T) {
    /* Step 1. Create new code */
    createReqBodyReader := strings.NewReader(`{"urlOriginal": "https://chopikus.github.io"}`)
    req := httptest.NewRequest(http.MethodPost, "/api/new")
    w := httptest.NewRecorder()


}
