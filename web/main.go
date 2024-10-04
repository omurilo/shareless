package web

import (
	"html/template"
	"net/http"
)

func Shareless(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/index.html")
	tmpl.Execute(w, nil)
}

func Share(w http.ResponseWriter, share interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/share.html")
	tmpl.Execute(w, share)
}

func Shared(w http.ResponseWriter, shared interface{}) {
	tmpl, _ := template.ParseFiles("web/templates/shared.html")
	tmpl.Execute(w, shared)
}
