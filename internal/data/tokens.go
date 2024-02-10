package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"

	"github.com/ivcp/polls/internal/validator"
)

type Token struct {
	Plaintext string
	Hash      []byte
}

func GenerateToken() (*Token, error) {
	var token Token

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return &token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
