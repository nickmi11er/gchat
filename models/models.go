package models

import (
	"encoding/json"
	"log"
)

// Channel model
type Channel struct {
}

// User model
type User struct {
	userID       string
	username     string
	passwordHash string
	channels     []Channel
}

type Response struct {
	Code   int               `json:"code"`
	Status string            `json:"status"`
	Params map[string]string `json:"params"`
}

func (r *Response) Json() []byte {
	m, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(m))
	return m
}
