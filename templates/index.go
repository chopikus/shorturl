package templates

import (
    "html/template"
)

var Index *template.Template

func init() {
    t, err := template.ParseFiles("templates/index.html")
    if err != nil {
        errorString := "index.html could not be parsed!" + err.Error()
        panic(errorString)
    }
    Index = t
}
