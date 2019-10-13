package model

import "time"

// Message ...
type Message struct {
	User *User
	Body string
	When time.Time
}
