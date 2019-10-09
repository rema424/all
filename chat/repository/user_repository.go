package repository

import (
	"chat/model"

	"github.com/jmoiron/sqlx"
)

// GetUserByID ...
func GetUserByID(db *sqlx.DB, userID int) (*model.User, error) {
	q := `
  SELECT id, full_name, username, mail
  FROM user
  WHERE id = ?;
  `
	var u model.User
	if err := db.Get(&u, q, userID); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByUserName ...
func GetUserByUserName(db *sqlx.Tx, username string) (*model.User, error) {
	q := `
  SELECT id, full_name, username, mail, password
  FROM user
  WHERE username = ?;
  `
	var u model.User
	if err := db.Get(&u, q, username); err != nil {
		return nil, err
	}
	return &u, nil
}
