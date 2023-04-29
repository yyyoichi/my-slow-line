package jwttoken

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtToken struct {
	secret string
	start  time.Time
}

func NewJwt(secret string) *JwtToken {
	return &JwtToken{secret: secret, start: time.Now()}
}

func (jt *JwtToken) Generate(userId string) (string, error) {
	rc := &jwt.RegisteredClaims{}
	rc.ExpiresAt = jwt.NewNumericDate(jt.start.AddDate(0, 0, 7))
	rc.ID = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	tokenString, err := token.SignedString([]byte(jt.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (jt *JwtToken) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jt.secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid")
	}

	rc, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid")
	}
	return rc, nil
}
