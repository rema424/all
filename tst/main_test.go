package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"tst/services/greeter"
)

var (
	server *httptest.Server
	client *http.Client
	gp     *greeter.Provider
	dbs    = make(map[string]map[string]interface{})
)

func TestMain(m *testing.M) {
	// データベース等の初期化
	e := newEcho()
	handler := routes(e)
	dbs["greeter"] = make(map[string]interface{})
	dbs["greeter"]["memory"] = greeter.NewMemoryDB()
	dbs["greeter"]["mysql"] = greeter.NewMysqlDB()

	// テストサーバーの起動
	server = httptest.NewServer(handler)
	defer server.Close()
	client = server.Client()

	os.Exit(m.Run())
}

func TestHandleGreet(t *testing.T) {
	for name, db := range dbs["greeter"] {
		t.Run(name, func(t *testing.T) {
			switch db := db.(type) {
			case *greeter.MemoryDB:
				gp.SetDB(db)
			case *greeter.MysqlDB:
				gp.SetDB(db)
			default:
				t.Fatal("invalid database")
			}
		})
	}
}

// func TestGreetRoutes(t *testing.T) {
// 	// データベース等の初期化
// 	e := newEcho()
// 	var db greeter.Database
// 	p := greeter.NewProvider(db)
// 	greetRoutes(e, p)

// 	// テストサーバーの起動
// 	s := httptest.NewServer(e)
// 	defer s.Close()
// 	c := s.Client()

// 	// リクエスト
// 	resp, err := c.Get(s.URL + "/greet?name=gopher")
// 	if err != nil {
// 		t.Fatalf("http.Get failed: %s", err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("resp.StatusCode: got: %d, want: %d", resp.StatusCode, http.StatusOK)
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	resp.Body.Close()
// 	if err != nil {
// 		t.Fatalf("ioutil.ReadAll failed: %s", err)
// 	}

// 	// アサーション
// 	got := string(body)
// 	want := "Hello, gopher!"
// 	if got != want {
// 		t.Fatalf("request: /greet?name=gopher, got %s, want %s", got, want)
// 	}
// }
