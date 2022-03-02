package main

import "github.com/twinj/uuid"

// generateCrypticState - generate a random string to re-evaluate spotify login
func generateCrypticState() string {
	id := uuid.NewV4()
	return id.String()
}
