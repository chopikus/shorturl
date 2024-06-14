package templates

import (
    "html/template"
)

var Redirect *template.Template

func init() {
    t, err := template.ParseFiles("templates/redirect.html")
    if err != nil {
        errorString := "redirect.html could not be parsed!" + err.Error()
        panic(errorString)
    }
    Redirect = t
}
