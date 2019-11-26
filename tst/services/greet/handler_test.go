package greet

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHandleGreet(t *testing.T) {
	// Setup
	e := echo.New()

	cases := []struct {
		in, out string
	}{
		{"alice", "Hello, alice!"},
		{"bob", "Hello, bob!"},
		{"carol", "Hello, carol!"},
		{"dave", "Hello, dave!"},
	}

	for i, cs := range cases {
		// Request
		req := httptest.NewRequest(http.MethodGet, "/greet?name="+cs.in, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := HandleGreet(c)

		// Assertions
		if err != nil {
			t.Fatalf("#%d: HandleGreet failed: %s", i, err)
		}
		if rec.Code != http.StatusOK {
			t.Fatalf("#%d: response code is not 200", i)
		}
		got, want := rec.Body.String(), cs.out
		if got != want {
			t.Fatalf("#%d: request: /greet?name=%s, got: %s, want: %s", i, cs.in, got, want)
		}
	}
}
