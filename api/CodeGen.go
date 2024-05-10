package main

import (
	"math/rand"
	"net/url"
	"time"
)

func GenerateVerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	code := string(b)
	encodedCode := url.QueryEscape(code)
	return encodedCode
}
