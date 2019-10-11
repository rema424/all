package service

import (
	"chat/infra"
	"chat/model"
	"chat/repository"
)

// GetAllRooms ...
func GetAllRooms(db *infra.DB) ([]*model.Room, error) {
	return repository.GetAllRooms(db)
}
