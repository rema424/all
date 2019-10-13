package model

// User はチャットを行う1人のユーザーを表します。
// 自分がどのルームへ所属しているかを認識しています。
type User struct {
	ID      int
	Name    string
	Email   string
	RoomIDs []int
}
