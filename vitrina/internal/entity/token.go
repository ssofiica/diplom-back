package entity

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

type JWT struct {
	Secret  []byte
	ExpTime time.Duration
}

func NewJWT(secret string, duration string) (JWT, error) {
	expiration, err := time.ParseDuration(duration)
	if err != nil {
		return JWT{
			Secret:  []byte(secret),
			ExpTime: time.Duration(48 * time.Hour),
		}, err
	}
	return JWT{
		Secret:  []byte(secret),
		ExpTime: expiration,
	}, nil
}

func (j JWT) GenerateToken(userID uint32, name, email string) (string, error) {
	claims := Claims{
		User: User{ID: userID, Name: name, Email: email},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.Secret)
}
