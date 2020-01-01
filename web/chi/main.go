/*
  mysql.server start
  mysql -uroot -e 'create database if not exists chi_echo_example;'
  mysql -uroot -e 'create user if not exists devuser@localhost identified by "Passw0rd!";'
  mysql -uroot -e 'grant all privileges on chi_echo_example.* to devuser@localhost;'
  mysql -uroot -e 'show databases;'
  mysql -uroot -e 'select host, user from mysql.user;'
  mysql -uroot -e 'show grants for devuser@localhost;'
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/flosch/pongo2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
)

const (
	dbName = "chi_echo_example"
	dbPass = "Passw0rd!"
	dbUser = "devuser"
	dbHost = "127.0.0.1"
	dbPort = "3306"
)

var (
	renderer  = NewRenderer("web/template")
	validate  = validator.New()
	db        *sqlx.DB
	shemaUser = `
create table if not exists user (
  id bigint auto_increment,
  email varchar(255) character set latin1 collate latin1_bin,
  password varchar(255),
  primary key (id),
  unique key (email)
);
`
	schemaSession = `
create table if not exists session (
  id varchar(255) character set latin1 collate latin1_bin,
  csrf varchar(255),
  user_id bigint,
  expire_at bigint default 0,
  primary key (id),
  foreign key (user_id) references user (id) on delete cascade on update cascade,
  key (user_id)
);
`
)

func main() {
	// DB
	cfg := mysql.Config{
		User:                 dbUser,
		Passwd:               dbPass,
		Net:                  "tcp",
		Addr:                 dbHost + ":" + dbPort,
		DBName:               dbName,
		Collation:            "utf8mb4_bin",
		InterpolateParams:    true,
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	db = DB(cfg)
	db.MustExec(shemaUser)
	db.MustExec(schemaSession)

	var eg errgroup.Group
	chiPort := "3000"
	// echoPort := "3001"

	// Chi
	eg.Go(func() error {
		r := Chi()
		log.Printf("chi: Listening on port %s", chiPort)
		return http.ListenAndServe(fmt.Sprintf(":%s", chiPort), r)
	})

	if err := eg.Wait(); err != nil {
		log.Fatalln(err)
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
// DB
// ------------------------------

func DB(cfg mysql.Config) *sqlx.DB {
	db, err := sqlx.Connect("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(60 * time.Second)
	return db
}

// ------------------------------
// Chi
// ------------------------------

func Chi() *chi.Mux {
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
	return r
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
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK) // 200

	w.Header().Add("Content-type", "text/html")
	renderer.Execute(w, "auth/signup.html", nil)
}

func HandlePostSignup(w http.ResponseWriter, r *http.Request) {
	// HTTPリクエストからパラメータを取得する。
	// バリデーションエラーが発生した場合はサインアップ画面を再表示する。
	in := SignupForm{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))

	// pp.Println(r)
	if err := validate.Struct(in); err != nil {
		log.Println(parseSignupFormError(err))
		w.WriteHeader(http.StatusUnprocessableEntity) // 422
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "auth/signup.html", map[string]interface{}{
			"Errors": parseSignupFormError(err),
		})
		return
	}

	// サインアップロジックを実行する。
	// サーバーエラーが発生した場合はエラーページを表示する。
	if err := doSomething(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError) // 500
		w.Header().Add("Content-type", "text/html")
		renderer.Execute(w, "error.html", nil)
		return
	}

	// サインアップに成功したらCookieにセッションIDを保存しつつ、
	// TOPページへリダイレクトする。
	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "<embed-session-id-here>",
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther) // 303
	return
}

func doSomething() error { return nil }
