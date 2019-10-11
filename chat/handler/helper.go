package handler

import (
	"chat/infra"
	"chat/model"
	"chat/service"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Room.ID: *model.Hub
var hubstore = make(map[int]*model.Hub)

var dbx *sqlx.DB

// Init ...
func Init(db *sqlx.DB) {
	dbx = db

	ctx := context.Background()
	mydb := infra.NewDB(ctx, dbx)
	rooms, err := service.GetAllRooms(mydb)
	fmt.Println(rooms, err)
}

// func GetExDB()
