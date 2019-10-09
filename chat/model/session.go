package model

import "time"

// Session ...
type Session struct {
	ID        string    `db:"id"`
	UserID    int       `db:"user_id"`
	LastLogin time.Time `db:"last_login"`
}
