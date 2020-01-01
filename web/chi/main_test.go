package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandlePostSignup(t *testing.T) {
	{
		// ----------
		// 422
		// ----------
		tests := []struct {
			email    string
			password string
			errors   []string
		}{
			{"", "", []string{"メールアドレスの入力は必須です", "パスワードの入力は必須です"}},
			{"aaaa", "aaaaabb", []string{"メールアドレスの形式が不正です", "パスワードは8文字以上15文字以内で入力してください"}},
			{"aaaaa@example.com", "aaaaabbbbbcccccd", []string{"パスワードは8文字以上15文字以内で入力してください"}},
		}
		for i, tt := range tests {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(form.Encode()))
			req.Header.Add("Content-type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			HandlePostSignup(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("#%d: handler returned wrong status code: got=%d, want=%d", i, rec.Code, http.StatusUnprocessableEntity)
			}

			body := rec.Body.String()
			for _, msg := range tt.errors {
				if !strings.Contains(body, msg) {
					t.Errorf("#d: response body does not contain '%s'", msg)
				}
			}
		}
	}

	{
		// ----------
		// 303
		// ----------
		tests := []struct {
			email    string
			password string
		}{
			{"aaa@example.com", "Passw0rd!"},
		}
		for i, tt := range tests {
			form := url.Values{}
			form.Add("email", tt.email)
			form.Add("password", tt.password)
			req := httptest.NewRequest(http.MethodPost, "/signup", strings.NewReader(form.Encode()))
			req.Header.Add("Content-type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			HandlePostSignup(rec, req)

			if rec.Code != http.StatusSeeOther {
				t.Errorf("#%d: handler returned wrong status code: got=%d, want=%d", i, rec.Code, http.StatusSeeOther)
			}

			resp := rec.Result()

			if u, err := resp.Location(); err != nil {
				t.Errorf("#%d: resp.Location() returned error: %s", i, err.Error())
			} else if u.Path != "/" {
				t.Errorf("#%d: wrong redirect url path. got=%s, want=%s", i, u.Path, "/")
			}

			cookies := &http.Request{Header: http.Header{"Cookie": rec.Header()["Set-Cookie"]}}
			if c, err := cookies.Cookie("session-id"); err != nil {
				t.Errorf("#%d: response does not contain cookie[%s]", i, "session-id")
			} else if c.Value == "" {
				t.Errorf("#%d: cookie value is empty", i)
			}
		}
	}
}
