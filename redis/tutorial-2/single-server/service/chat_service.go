// クライアントがルームに接続を試みる
// ルームがアクティブか判定する（非アクティブなら起動する）
// クライアントはルームに参加する
// クライアントは出版を待機する
// クライアントは購読を待機する
// クライアントはルームに対してメッセージを発行する（クライアントの出版）
// ルームは接続中の全クライアントを確認する
// ルームは接続中の全クライアントにメッセージを送信する（クライアントの購読）

package service

import (
	"fmt"
	"single-server/infra"
	"single-server/model"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const messageBufferSize = 256

type client struct {
	socket *websocket.Conn
	send   chan *model.Message
	// room   *model.Room
	// user   *model.User
	roomSpvr *roomSpvr
	userID   int
}

type roomSpvr struct {
	roomID int
	sync.Once
	sync.RWMutex
	writeCh chan *model.Message
	joinCh  chan *client
	leaveCh chan *client
	doneCh  chan struct{}
	clients map[*client]bool
}

var gRoomSpvrCache = &struct {
	sync.RWMutex
	roomSpvrs map[int]*roomSpvr
}{
	RWMutex:   sync.RWMutex{},
	roomSpvrs: make(map[int]*roomSpvr),
}

var roomSpvrPool = &sync.Pool{
	New: func() interface{} {
		return &roomSpvr{
			roomID:  0,
			Once:    sync.Once{},
			RWMutex: sync.RWMutex{},
			writeCh: make(chan *model.Message),
			joinCh:  make(chan *client),
			leaveCh: make(chan *client),
			clients: make(map[*client]bool),
			doneCh:  make(chan struct{}),
		}
	},
}

func newRoomSpvr(roomID int) *roomSpvr {
	r := roomSpvrPool.New().(*roomSpvr)
	r.roomID = roomID
	return r
}

func releaseRoomSpvr(r *roomSpvr) {
	r.Lock()
	defer r.Unlock()

	close(r.doneCh)
	close(r.writeCh)
	close(r.joinCh)
	close(r.leaveCh)

	r.roomID = 0
	r.Once = sync.Once{}
	r.doneCh = make(chan struct{})
	r.writeCh = make(chan *model.Message)
	r.joinCh = make(chan *client)
	r.leaveCh = make(chan *client)
	r.clients = make(map[*client]bool)

	roomSpvrPool.Put(r)
}

// ConnectChatRoom ...
func ConnectChatRoom(db *infra.DB, socket *websocket.Conn, roomID, userID int) {
	// チャットルームに接続する
	// ユーザーが入ろうとしている部屋へ参加資格があるかを確認する
	roomSpvr := getRoomSpvr(roomID)
	roomSpvr.Once.Do(func() { go roomSpvr.run() })

	client := &client{
		socket:   socket,
		send:     make(chan *model.Message, model.MessageBufferSize),
		roomSpvr: roomSpvr,
		userID:   userID,
	}

	roomSpvr.joinCh <- client
	defer func() { roomSpvr.leaveCh <- client }()

	go client.publish()
	client.subscribe()
}

func getRoomSpvr(roomID int) *roomSpvr {
	gRoomSpvrCache.Lock()
	defer gRoomSpvrCache.Unlock()

	var r *roomSpvr

	if val, ok := gRoomSpvrCache.roomSpvrs[roomID]; ok {
		fmt.Println("roomID:", roomID, "- roomをキャッシュから取得しました。")
		r = val
	} else {
		fmt.Println("roomID:", roomID, "- roomを新規に作成します。")
		r = newRoomSpvr(roomID)
		gRoomSpvrCache.roomSpvrs[roomID] = r
	}

	return r
}

func (r *roomSpvr) run() {
	fmt.Println("roomID:", r.roomID, "- run()を開始します。")
	defer fmt.Println("roomID:", r.roomID, "- run()を終了します。")

	go func() {
	L:
		for {
			select {
			case <-time.Tick(10 * time.Second):
				fmt.Println("roomID:", r.roomID, "- roomが起動中です。")
				if len(r.clients) == 0 {
					fmt.Println("roomID:", r.roomID, "- 参加者が0人です。部屋の監視を停止します。")
					for client := range r.clients {
						close(client.send)
					}
					close(r.doneCh)
					delete(gRoomSpvrCache.roomSpvrs, r.roomID)
					break L
				}
			}
		}
	}()

L:
	for {
		select {
		case client := <-r.joinCh:
			r.clients[client] = true
			fmt.Println("roomID:", r.roomID, "- userID:", client.userID, "がチャットに参加しました")
			r.describe()
		case client := <-r.leaveCh:
			delete(r.clients, client)
			close(client.send)
			fmt.Println("roomID:", r.roomID, "- userID:", client.userID, "が退室しました")
			r.describe()
		case msg := <-r.writeCh:
			fmt.Printf("メッセージを受信しました\n")
			for client := range r.clients {
				select {
				case client.send <- msg:
					fmt.Printf(" -- クライアントに送信しました\n")
				default:
					delete(r.clients, client)
					close(client.send)
					fmt.Printf(" -- 送信に失敗しました。クライアントをクリーンアップします。\n")
				}
			}
		case <-r.doneCh:
			break L
		}
	}
}

func (r *roomSpvr) describe() {
	r.RLock()
	defer r.RUnlock()

	userIDs := make([]string, 0, len(r.clients))
	for client := range r.clients {
		userIDs = append(userIDs, strconv.Itoa(client.userID))
	}

	fmt.Printf("Room: %d\n", r.roomID)
	fmt.Printf("--- 参加人数: %d\n", len(userIDs))
	fmt.Printf("--- 参加者: %s\n", strings.Join(userIDs, ", "))
}

func (c *client) subscribe() {
	fmt.Println("userID:", c.userID, "- 購読を開始しました。")
	defer c.socket.Close()
	defer fmt.Println("userID:", c.userID, "- 購読を終了しました。")

	for {
		var msg *model.Message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			c.roomSpvr.writeCh <- msg
		} else {
			break
		}
	}
}

func (c *client) publish() {
	fmt.Println("userID:", c.userID, "- 出版を開始しました。")
	defer c.socket.Close()
	defer fmt.Println("userID:", c.userID, "- 出版を終了しました。")

	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
}

// // ConnectClient ...
// func ConnectClient(db *infra.DB, socket *websocket.Conn, roomID, userID int) {
// 	user, err := GetUserByID(db, userID)
// 	if err != nil {
// 		fmt.Println("Userの取得に失敗しました：", user, "-", err)
// 	}

// 	roomSpvr, err := NewRoomManager(db, roomID, userID)

// 	if err != nil {
// 		fmt.Println("Roomの取得に失敗しました：", room, "-", err)
// 	}

// 	client := &clientSupervisor{
// 		socket: socket,
// 		send:   make(chan *model.Message, model.MessageBufferSize),
// 		room:   room,
// 		user:   user,
// 	}

// 	room.JoinCh <- client
// 	defer func() { room.LeaveCh <- client }()
// 	go client.Write()
// 	client.Read()

// }

// func (c *clientSupervisor) Read() {
// 	defer c.socket.Close()

// 	fmt.Println("client read started")
// 	defer fmt.Println("client read finished")

// 	for {
// 		fmt.Println("client read loop")
// 		var msg *model.Message
// 		if err := c.socket.ReadJSON(&msg); err == nil {
// 			msg.When = time.Now()
// 			msg.User = c.user
// 			c.room.MessageCh <- msg
// 		} else {
// 			break
// 		}
// 	}
// }

// func (c *clientSupervisor) Write() {
// 	defer c.socket.Close()

// 	fmt.Println("client write started")
// 	defer fmt.Println("client write finished")

// 	for msg := range c.send {
// 		fmt.Println("client write loop")
// 		if err := c.socket.WriteJSON(msg); err != nil {
// 			break
// 		}
// 	}
// }

// // newRoomSpvr ...
// func newRoomSpvr(db *infra.DB, roomID, userID int) *roomSpvr {
// 	return nil
// }

// // GetRoomByID ...
// func GetRoomByID(db *infra.DB, roomID int) (*model.Room, error) {
// 	fmt.Println("TODO implement service.GetRoomByID")
// 	return nil, nil
// }

// // JoinRoom ...
// // func JoinRoom(db *infra.DB, roomID int, socket *websocket.Conn, user *model.User) {
// // 	room, err := GetRoomByID(db, roomID)
// // 	if err != nil {
// // 		fmt.Println("Roomの取得に失敗しました：", room, "-", err)
// // 	}

// // 	client := &model.Client{
// // 		Socket: socket,
// // 		Send:   make(chan *model.Message, model.MessageBufferSize),
// // 		Room:   room,
// // 		User:   user,
// // 	}

// // 	room.JoinCh <- client
// // 	defer func() { room.LeaveCh <- client }()
// // 	go client.Write()
// // 	client.Read()
// // }

// // GetRoomByID ...
// // func GetRoomByID(db *infra.DB, roomID int) (*model.Room, error) {
// // 	return nil, nil
// // }
