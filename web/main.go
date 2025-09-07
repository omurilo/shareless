package web

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed templates/*.html
var templates embed.FS

func Shareless(w http.ResponseWriter, r *http.Request) {
	tmpl, _:= template.ParseFS(templates, "templates/index.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Shared(w http.ResponseWriter, shared interface{}) {
	tmpl, _ := template.ParseFS(templates, "templates/shared.html")
	err := tmpl.Execute(w, shared)
	if err != nil {
		log.Fatalln(err)
	}
}

func About(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFS(templates, "templates/about.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func Privacy(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFS(templates, "templates/privacy.html")
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
