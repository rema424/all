package repository

import (
	"chat/model"

	"github.com/jmoiron/sqlx"
)

// GetSessionByID ...
func GetSessionByID(db *sqlx.DB, sessionID int) (*model.Session, error) {
	q := `
  SELECT id, user_id, last_login
  FROM session
  WHERE id = ?;
  `
	var s model.Session
	if err := db.Get(&s, q, sessionID); err != nil {
		return nil, err
	}
	return &s, nil
}
