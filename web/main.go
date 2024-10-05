package web

import (
	"html/template"
	"net/http"
)

func Shareless(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/index.html")
	tmpl.Execute(w, nil)
}

func Shared(w http.ResponseWriter, shared interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/shared.html")
	tmpl.Execute(w, shared)
}

func About(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/about.html")
	tmpl.Execute(w, nil)
}
