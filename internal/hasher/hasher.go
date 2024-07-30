package hasher

import (
	"crypto/sha512"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(todoPassword string) (string, error) {
	secret := []byte("secret_key")
	pass := []byte(todoPassword)
	hasher := sha512.New()
	pass = append(pass, []byte("salt")...)
	hash := hasher.Sum(pass)
	hashString := hex.EncodeToString(hash)

	claims := jwt.MapClaims{
		"hash": hashString,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
