package main

import (
    "os"
    "log"
    "time"
    "math/rand"
    "database/sql"
    _ "github.com/lib/pq"
)

var db *sql.DB

type Short struct {
    UrlOriginal string `json:"urlOriginal"`
    UrlCode string `json:"urlCode"`
    ExpiresOn time.Time `json:"expiresOn"`
}

func init() {
    /* rand is used to generate new URL codes */
    rand.Seed(time.Now().UnixNano())
    
    user := os.Getenv("SHORTURL_POSTGRES_USER")
    password := os.Getenv("SHORTURL_POSTGRES_PASSWORD")
    host := os.Getenv("SHORTURL_POSTGRES_HOST")
    if host == "" {
        host = "/var/run/postgresql"
    }

    openString := fmt.Sprintf("dbname=urldb user=%s password=%s host=%s sslmode=disable", user, password, host)
    tmpDB, err := sql.Open("postgres", openString)
    if err != nil {
        log.Fatal(err)
    }
    
    // testing the connection, sql.Open doesn't catch the errors
    _, err = tmpDB.Exec("SELECT 2+3")
    if err != nil {
        log.Fatal(err)
    }

    db = tmpDB
}

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
func createCode(url string) (Short, error) {
    var code string = generateCode()
    var exp time.Time

    err := db.
           QueryRow(`INSERT INTO urls 
                     (url_original, url_code, expires_on)
                     VALUES ($1, $2, NOW() + '24 hour')
                     RETURNING expires_on`, 
                    url, code).
           Scan(&exp)
            
    
    if err != nil {
        return Short{}, err
    }
    
    return Short{UrlCode: code, UrlOriginal: url, ExpiresOn: exp}, nil
}


// Gets a short from the database for a provided code
func getCode(code string) (Short, error) {
    var url string
    var exp time.Time

    err := db.
           QueryRow("SELECT url_original, expires_on FROM urls WHERE expires_on >= NOW() AND url_code=$1", code).
           Scan(&url, &exp)

    if err != nil {
        return Short{}, err
    }

    return Short{UrlCode: code, UrlOriginal: url, ExpiresOn: exp}, nil
}

func removeExpiredCodes() error {
    _, err := db.Exec("DELETE FROM urls WHERE expires_on < NOW()")
    if err != nil {
        return err
    }
    return nil
}
