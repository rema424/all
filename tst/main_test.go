package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoute(t *testing.T) {
	mux := CreateDefaultMux()
	s := httptest.NewServer(Route(mux))
	defer s.Close()

	res, err := http.Get(s.URL + "/greet?name=gopher")
	if err != nil {
		t.Fatalf("http.Get failed: %s", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("res.StatusCode: got: %d, want: %d", res.StatusCode, http.StatusOK)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatalf("ioutil.ReadAll failed: %s", err)
	}

	got := string(body)
	want := "Hello, gopher!"
	if got != want {
		t.Fatalf("request: /greet?name=gopher, got %s, want %s", got, want)
	}
}
