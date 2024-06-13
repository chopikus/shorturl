package main

import (
    "math/rand"
)

func generateCode() string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
    const keyLength = 6

    shortKey := make([]byte, keyLength)
    for i := range shortKey {
        shortKey[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortKey)
}

// Generates a new code for the provided url and puts it into the database.
func createCode(url string) (string, error) {
    var code string = generateCode()

    _, err := db.Exec("INSERT INTO urls (url_original, url_code) VALUES ($1, $2)", url, code)
    
    if err != nil {
        return "", err
    }

    return code, nil
}


// Gets a url from the database for a provided code
func getUrl(code string) (string, error) {
    var url string

    err := db.
           QueryRow("SELECT url_original FROM urls WHERE url_code=$1", code).
           Scan(&url)

    if err != nil {
        return "", err
    }

    return url, nil
}