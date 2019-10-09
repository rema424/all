package model

import "time"

// User ...
type User struct {
	ID       int        `json:"id"`
	Fullname string     `json:"fullname"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Password string     `json:"-"`
	Role     string     `json:"role"`
	MuteDate *time.Time `json:"mute_date"`
}
