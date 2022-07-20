package cache

import (
	"crypto/sha1"
	"fmt"
)

// Sha1 Calculate the sha1 hash of a string
func Sha1(str string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(str))
	return fmt.Sprintf("%x", h.Sum(nil))
}
