package token

import (
	"crypto/rand"
	"fmt"
)

func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)

	return fmt.Sprintf("%x", b)
}
