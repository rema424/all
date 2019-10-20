package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// Session ...
type Session struct {
	SessionID      string    `db:"session_id"`
	UserID         int       `db:"user_id"`
	LastLoggedInAt time.Time `db:"last_logged_in_at"`
}

// --------------------------------------------------
// signupExecHandler
// --------------------------------------------------

type signupExecIn struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	PasswordConfirm string `form:"password-confirm"`
}

type signupExecOut struct {
	Message string `json:"message"`
}

func signupExecHandler(c echo.Context) error {
	in := signupExecIn{}
	out := signupExecOut{}

	if err := c.Bind(&in); err != nil {
		fmt.Println(err.Error())
		out.Message = "不正なリクエストです。"
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&in); err != nil {
		fmt.Println(err.Error())
		out.Message = "不正な値です。"
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	if in.Password != in.PasswordConfirm {
		out.Message = "パスワードが一致しません。"
		return c.JSON(http.StatusForbidden, out)
	}

	dbx := GetDBx(c)

	q := `select email from user where email = ?;`
	var email string
	if err := dbx.Get(&email, q, in.Email); err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err.Error())
			out.Message = "そのメールアドレスは既に使われています。"
			return c.JSON(http.StatusForbidden, out)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		out.Message = "予期せぬエラーが発生しました。"
		return c.JSON(http.StatusInternalServerError, out)
	}

	fmt.Println(in.Email)
	u := User{
		Name:         namesgenerator.GetRandomName(0),
		Email:        in.Email,
		PasswordHash: string(hash),
	}
	var sessID string

	err = dbx.BeginTx(func(x *DBx) error {
		q := `insert into user (name, email, password_hash) values (:name, :email, :password_hash);`
		res, err := x.NamedExec(q, u)
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		u.ID = int(id)
		sessID = uuid.New().String()

		q = `insert into session (session_id, user_id) values (?, ?);`
		_, err = x.Exec(q, sessID, u.ID)
		return err
	})

	if err != nil {
		fmt.Println(err.Error())
		out.Message = "予期せぬエラーが発生しました。"
		return c.JSON(http.StatusInternalServerError, out)
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessid",
		Value:    sessID,
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // 開発環境ではfalse
	})

	out.Message = "サインアップに成功しました。"
	return c.JSON(http.StatusOK, out)
}

// --------------------------------------------------
// loginExecHandler
// --------------------------------------------------

type loginExecIn struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type loginExecOut struct {
	Message string `json:"message"`
}

func loginExecHandler(c echo.Context) error {
	in := loginExecIn{}
	out := loginExecOut{}

	if err := c.Bind(&in); err != nil {
		fmt.Println(err.Error())
		out.Message = "メールアドレスまたはパスワードが違います。"
		return c.JSON(http.StatusBadRequest, out)
	}

	if err := c.Validate(&in); err != nil {
		fmt.Println(err.Error())
		out.Message = "メールアドレスまたはパスワードが違います。"
		return c.JSON(http.StatusUnprocessableEntity, out)
	}

	dbx := GetDBx(c)

	q := `
  select id, coalesce(password_hash, '') as password_hash
  from user
  where email = ?;`
	var u User
	if err := dbx.Get(&u, q, in.Email); err != nil {
		fmt.Println(err.Error())
		out.Message = "メールアドレスまたはパスワードが違います。"
		return c.JSON(http.StatusForbidden, out)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		fmt.Println(err.Error())
		out.Message = "メールアドレスまたはパスワードが違います。"
		return c.JSON(http.StatusForbidden, out)
	}

	sessID := uuid.New().String()
	err := dbx.BeginTx(func(x *DBx) error {
		q := `insert into session (session_id, user_id) values (?, ?);`
		_, err := x.Exec(q, sessID, u.ID)
		return err
	})

	if err != nil {
		fmt.Println(err.Error())
		out.Message = "予期しないエラーが発生しました。"
		return c.JSON(http.StatusInternalServerError, out)
	}

	c.SetCookie(&http.Cookie{
		Name:     "sessid",
		Value:    sessID,
		Expires:  time.Now().Add(60 * 24 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		// Secure:   true, // 開発環境ではfalse
	})

	out.Message = "ログインに成功しました。"
	return c.JSON(http.StatusOK, out)
}

// --------------------------------------------------
// isLoggedInHandler
// --------------------------------------------------

func isLoggedInHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessid")
		if err != nil {
			fmt.Println("Cookieが見つかりません。")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		dbx := GetDBx(c)

		sessID := cookie.Value
		var s Session
		q := `
    select session_id, user_id, last_logged_in_at
    from session
    where session_id = ?;`
		if err := dbx.Get(&s, q, sessID); err != nil {
			fmt.Println("セッションが見つかりません。")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if time.Now().After(s.LastLoggedInAt.AddDate(0, 0, 60)) {
			fmt.Println("セッションの有効期限が切れています。")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		if time.Now().After(s.LastLoggedInAt.Add(1 * time.Hour)) {
			fmt.Println("最終ログインから1時間経過しています。最終ログイン時刻を更新します。")
			err := dbx.BeginTx(func(x *DBx) error {
				q := `update session set last_logged_in_at = current_timestamp() where session_id = ?;`
				_, err := x.Exec(q, sessID)
				return err
			})
			if err != nil {
				fmt.Println("予期せぬエラーが発生しました。")
				return c.Redirect(http.StatusSeeOther, "/login")
			}
		}

		return next(c)
	}
}
