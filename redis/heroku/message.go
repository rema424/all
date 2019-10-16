package main

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type message struct {
	Handle string `json:"handle"`
	Text   string `json:"text"`
}

func validateMessage(data []byte) (message, error) {
	var msg message
	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, errors.Wrap(err, "Unmarshaling message")
	}

	if msg.Handle == "" && msg.Text == "" {
		return msg, errors.New("Message has no Handle or Text")
	}

	return msg, nil
}
