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
	email           string `form:"email"`
	password        string `form:"password"`
	passwordConfirm string `form:"passwordConfirm"`
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

	if in.password != in.passwordConfirm {
		out.Message = "パスワードが一致しません。"
		return c.JSON(http.StatusForbidden, out)
	}

	dbx := GetDBx(c)

	q := `select email from user where email = ?;`
	var email string
	if err := dbx.Get(&email, q, in.email); err != nil {
		if err != sql.ErrNoRows {
			fmt.Println(err.Error())
			out.Message = "そのメールアドレスは既に使われています。"
			return c.JSON(http.StatusForbidden, out)
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err.Error())
		out.Message = "予期せぬエラーが発生しました。"
		return c.JSON(http.StatusInternalServerError, out)
	}

	u := User{
		Name:         namesgenerator.GetRandomName(0),
		Email:        in.email,
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
	email    string `form:"email"`
	password string `form:"password"`
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
	if err := dbx.Get(&u, q, in.email); err != nil {
		fmt.Println(err.Error())
		out.Message = "メールアドレスまたはパスワードが違います。"
		return c.JSON(http.StatusForbidden, out)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.password)); err != nil {
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
