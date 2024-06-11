package main

import (
  "database/sql"
  "log"
  _ "github.com/lib/pq"
  "fmt"
  "math/rand"
  "time"
)

var db *sql.DB

/* Short represents shortening of the link. */
type Short struct {
    urlCode string
    urlOriginal string
}

func generateShortKey() string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
    const keyLength = 6

    shortKey := make([]byte, keyLength)
    for i := range shortKey {
        shortKey[i] = charset[rand.Intn(len(charset))]
    }
    return string(shortKey)
}


func createShort(urlOriginal string) Short {
    var urlCode string = generateShortKey()

    _, err := db.Exec("INSERT INTO urls (url_original, url_code) VALUES ($1, $2)", urlOriginal, urlCode)
    
    if err != nil {
        log.Fatal(err)
    }

    return Short{urlCode: urlCode, urlOriginal: urlOriginal}
}

func getShort(urlCode string) Short {
    var urlOriginal string

    err := db.
           QueryRow("SELECT url_original FROM urls WHERE url_code=$1", urlCode).
           Scan(&urlOriginal)

    if err != nil {
        log.Fatal(err)
    }

    return Short{urlCode: urlCode, urlOriginal: urlOriginal}

}

func main() {
    rand.Seed(time.Now().UnixNano())
    tmpDB, err := sql.Open("postgres", "dbname=urldb user=ihor password=ihor host=/tmp sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    db = tmpDB
    //fmt.Println(createShort("https://linkedin.com"))
    fmt.Println(getShort("GML2JA"))
}
