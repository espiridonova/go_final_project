package api

import (
	"crypto/sha256"
	"fmt"
)

const secretKey = "lol123"

func createHash(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}
