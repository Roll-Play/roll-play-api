package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func CreateJWT(id uuid.UUID, expire_at time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "rollplay",
		"sub": id.String(),
		"exp": time.Now().Add(expire_at * time.Millisecond).Unix(),
	})

	key := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(key))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}
