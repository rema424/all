package model

// Room はチャットルーム1部屋を表します。
type Room struct {
	MessageCh chan *Message
	JoinCh    chan *Client
	LeaveCh   chan *Client
	Clients   map[*Clients]bool
}
