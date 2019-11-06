package main

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("unable to generate 16 bytes of randomness")
	}
	return hex.EncodeToString(b)
}
