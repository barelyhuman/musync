package main

import "github.com/twinj/uuid"

// GenerateCrypticState - generate a random string to re-evaluate spotify login
func GenerateCrypticState() string {
	id := uuid.NewV4()
	return id.String()
}
