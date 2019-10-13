package service

import (
	"fmt"
	"single-server/infra"
	"single-server/model"
)

// GetRoomByID ...
func GetRoomByID(db *infra.DB, roomID int) (*model.Room, error) {
	fmt.Println("TODO implement service.GetRoomByID")
	return nil, nil
}
