package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct{}

const secret = "test"

var testClaims = jwt.MapClaims{
	"sub": "102",
	"exp": time.Now().Add(time.Hour).Unix(),
	"iat": time.Now().Unix(),
	"nbf": time.Now().Unix(),
	"iss": "test-aud",
	"aud": "test-aud",
}

func (ta *TestAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ta *TestAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
}
