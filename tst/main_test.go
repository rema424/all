package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"tst/services/greeter"
)

func TestRoutes(t *testing.T) {
	// データベース等の初期化
	e := newEcho()
	handler := routes(e)

	// テストサーバーの起動
	s := httptest.NewServer(handler)
	defer s.Close()
}

func TestGreetRoutes(t *testing.T) {
	// データベース等の初期化
	e := newEcho()
	var db greeter.Database
	p := greeter.NewProvider(db)
	greetRoutes(e, p)

	// テストサーバーの起動
	s := httptest.NewServer(e)
	defer s.Close()
	c := s.Client()

	// リクエスト
	resp, err := c.Get(s.URL + "/greet?name=gopher")
	if err != nil {
		t.Fatalf("http.Get failed: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("resp.StatusCode: got: %d, want: %d", resp.StatusCode, http.StatusOK)
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("ioutil.ReadAll failed: %s", err)
	}

	// アサーション
	got := string(body)
	want := "Hello, gopher!"
	if got != want {
		t.Fatalf("request: /greet?name=gopher, got %s, want %s", got, want)
	}
}
