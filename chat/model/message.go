package model

import "time"

// Message ...
type Message struct {
	ID        int       `json:"id"`
	Action    string    `json:"action"`
	Sender    *User     `json:"sender"`
	Recipient *User     `json:"recipient"`
	Text      string    `json:"text"`
	SendDate  time.Time `json:"send_date"`
	Room      *Room     `json:"-"`
}
