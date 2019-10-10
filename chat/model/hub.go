package model

// Hub ...
type Hub struct {
	clients    map[*WSClient]bool
	room       *Room
	message    chan *Message
	ctl        chan *Message
	register   chan *Register
	unregister chan *Unregister
}

// Register ...
type Register struct {
	client *WSClient
}

// Unregister ...
type Unregister struct {
	client *WSClient
	msg    string
}
