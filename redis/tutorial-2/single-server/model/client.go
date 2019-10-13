package model

import "github.com/gorilla/websocket"

// Client はブラウザを初めとする端末を表します。
// WebSocket のコネクションを持ちます。
// 1人の User は複数の Client を持ちます。
type Client struct {
	Socket *websocket.Conn
	Send   chan *Message
	Room   *Room
	User   *User
}
