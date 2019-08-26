package passwordhash

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type HashSha256 struct{}

func NewSha256() *HashSha256 {
	return &HashSha256{}
}

func (h *HashSha256) Hash(s ...string) string {
	hash := hmac.New(sha256.New, []byte{})
	for _, v := range s {
		hash.Write([]byte(v))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (h *HashSha256) Salt() string {
	salt := make([]byte, DefaultPasswordSaltSize)
	rand.Read(salt)
	return hex.EncodeToString(salt)
}
