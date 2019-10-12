package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		fmt.Println(authCookie.Value, objx.MustFromBase64(authCookie.Value))
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

func main() {
	addr := flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	// Gomniauithのセットアップ
	gomniauth.SetSecurityKey(os.Getenv("GOMNIAUTH_SECURITY_KEY"))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"), "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// チャットルームを開始します
	go r.run()

	// Webサーバーを起動します
	log.Println("Webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
