package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/k0kubun/pp"
)

var (
	renderer = NewRenderer("web/template")
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "top.html", nil)
	})

	r.Get("/signup", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "auth/signup.html", nil)
	})

	r.Post("/signup", func(w http.ResponseWriter, r *http.Request) {
		// pp.Println(r)
		in := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		pp.Println(in)
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "auth/signup.html", nil)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "auth/login.html", nil)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "auth/login.html", nil)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "ok",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatal(err)
	}
}

// ------------------------------
// Renderer
// ------------------------------

func NewRenderer(tmplDir string) *Renderer {
	return &Renderer{tmplDir}
}

type Renderer struct {
	TmplDir string
}

func (r *Renderer) Execute(w io.Writer, name string, data map[string]interface{}) error {
	b, err := pongo2.Must(pongo2.FromCache(filepath.Join(r.TmplDir, name))).ExecuteBytes(data)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

// ------------------------------
// SignupForm
// ------------------------------

type SignupForm struct {
	Email    string `validate:"-"`
	Password string `validate:"-"`
}
