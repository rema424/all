package greeter

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

var (
	e   = echo.New()
	p   = NewProvider(nil)
	dbs = map[string]Database{}
)

func TestMain(m *testing.M) {
	dbs["memory"] = NewMemoryDB()
	dbs["mysql"] = NewMysqlDB()
	os.Exit(m.Run())
}

func TestHandleGreet(t *testing.T) {
	for name, db := range dbs {
		t.Run(name, func(t *testing.T) {
			p.db = db

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
				err := p.HandleGreet(c)

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
		})
	}
}
