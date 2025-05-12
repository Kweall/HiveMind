package resolver

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

// generateID: Генерация безопасного 16-символьного уникального ID
func GenerateID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("failed to generate secure ID: %v", err)
	}
	return hex.EncodeToString(b)
}
