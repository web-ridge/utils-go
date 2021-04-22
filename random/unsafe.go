package random

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededUnsafeRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func unsafeStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededUnsafeRand.Intn(len(charset))]
	}
	return string(b)
}

// UnsafeString returns a random string which can be used for things like file names. Don't use this for tokens or
// security since we should use match/crypto for that!
func UnsafeString(length int) string {
	return unsafeStringWithCharset(length, charset)
}
