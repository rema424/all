package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	tpl *template.Template
)

func main() {
	tpl = template.Must(template.ParseGlob("web/template/*"))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	http.ListenAndServe(":3333", r)
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
  tpl.
}
