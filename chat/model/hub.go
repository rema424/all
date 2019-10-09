package model

// Hub ...
type Hub struct {
	clients    map[*Client]bool
	room       *Room
	message    chan *Message
	ctl        chan *Message
	register   chan *Register
	unregister chan *Unregister
}

// Register ...
type Register struct {
	client *Client
}

// Unregister ...
type Unregister struct {
	client *Client
	msg    string
}
