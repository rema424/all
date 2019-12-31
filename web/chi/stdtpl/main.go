package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/k0kubun/pp"
	"github.com/go-chi/chi/middleware"
)

var (
	tpl *template.Template
)

func main() {
	tpl = template.Must(template.ParseGlob("template/*.html"))
	template.Must(tpl.ParseGlob("template/user/*.html"))
	// tpl = template.Must(template.ParseGlob("template/**/*.html"))
	// parseTemplate(nil)
	tpl.ParseFiles()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("template/base.html", "template/top.html"))
		// w.Write([]byte("hello world"))
    t.Execute(w, nil)
    pp.Println(t)
	})
	r.Get("/users", HandleUserList)
	http.ListenAndServe(":3333", r)
}

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "base.html", nil)
}

func HandleUserList(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "base", nil)
}

func parseTemplate(funcMap template.FuncMap) *template.Template {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	tplDir := filepath.Join(wd, "template")
	cleanRoot := filepath.Clean(tplDir)
	fmt.Println(wd, tplDir, cleanRoot)
	return nil
}
