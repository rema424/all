package service

import (
	"fmt"
	"single-server/infra"
	"single-server/model"
)

// GetUserByID ...
func GetUserByID(db *infra.DB, userID int) (*model.User, error) {
	fmt.Println("TODO implement service.GetUserByID")
	return nil, nil
}
