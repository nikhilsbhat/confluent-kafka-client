package main

import (
	"encoding/json"

	"github.com/google/uuid"
)

type fakeMessage struct {
	Message  string `json:"message"`
	Instance int    `json:"instance"`
}

func randStringBytes() string {
	uid := uuid.New()
	return uid.String()
}

func getBytes(index int, msg string) []byte {
	data := fakeMessage{Message: msg, Instance: index}
	message, _ := json.Marshal(&data)
	return message
}
