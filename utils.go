package main

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func PasswordHash(password string) string {
	h := crypto.SHA256.New()
	h.Write([]byte(password))
	bs := h.Sum([]byte{})
	return hex.EncodeToString(bs)
}
