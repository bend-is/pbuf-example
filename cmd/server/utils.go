package main

import (
	"github.com/gofrs/uuid"
	"log"
)

func generateNewID() (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		log.Printf("[ERROR] uuid generation error: %v", err)
		return "", err
	}

	return id.String(), nil
}
