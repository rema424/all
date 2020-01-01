package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/k0kubun/pp"
)

var (
	renderer = NewRenderer("web/template")
	validate = validator.New()
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", HandleGetTop)
	r.Get("/signup", HandleGetSignup)
	r.Post("/signup", HandlePostSignup)

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
// Validate
// ------------------------------

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

// ------------------------------
// Top
// ------------------------------

func HandleGetTop(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	renderer.Execute(w, "top.html", nil)
}

// ------------------------------
// Signup
// ------------------------------

type SignupForm struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=15"`
}

var SignupFormFields = map[string]string{
	"Email":    "メールアドレス",
	"Password": "パスワード",
}

func parseSignupFormError(err error) []ValidationError {
	errs := make([]ValidationError, 0, 2)
	for _, err := range err.(validator.ValidationErrors) {
		f := err.Field()
		t := err.Tag()
		var msg string
		switch t {
		case "required":
			msg = fmt.Sprintf("%sの入力は必須です。", SignupFormFields[f])
		case "email":
			msg = fmt.Sprintf("%sの形式が不正です。", SignupFormFields[f])
		default:
			switch f {
			case "Password":
				if t == "min" || t == "max" {
					msg = fmt.Sprintf("%sは8文字以上15文字以内で入力してください。", SignupFormFields[f])
				}
			}
		}
		if msg != "" {
			errs = append(errs, ValidationError{f, t, msg})
		}
	}
	return errs
}

func HandleGetSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text/html")
	renderer.Execute(w, "auth/signup.html", nil)
}

func HandlePostSignup(w http.ResponseWriter, r *http.Request) {
	in := SignupForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	if err := validate.Struct(in); err != nil {
		log.Println(parseSignupFormError(err))
	}

	pp.Println(in)
	w.Header().Add("Content-type", "text/html")
	renderer.Execute(w, "auth/signup.html", nil)
}
