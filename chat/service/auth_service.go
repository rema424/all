package service

import (
	"chat/repository"

	"github.com/jmoiron/sqlx"
)

// --------------------------------------------------
// Login
// --------------------------------------------------

// Login はログインを実行します。
func Login(db *sqlx.DB, username, password string) LoginOut {
	_, err := repository.GetUserByUserName(db, username)
	if err != nil {

	}

	return LoginOut{}
}

// LoginOut ...
type LoginOut struct {
}

// --------------------------------------------------
// Logout
// --------------------------------------------------

// Logout はログアウトを実行します。
func Logout() {

}

// --------------------------------------------------
// IsLoggedIn
// --------------------------------------------------

// IsLoggedIn は現在ログイン中か判定します。
func IsLoggedIn() {

}
