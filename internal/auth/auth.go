package auth

import (
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (hash string, err error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (match bool, err error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
