package app

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Session struct {
	UserId  string `json:"userId"`
	Channel string `json:"channel"`
	jwt.StandardClaims
}

func getSession(key string, secret string) (*Session, error) {
	token, err := jwt.ParseWithClaims(key, &Session{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return token.Claims.(*Session), nil
	}
	return nil, errors.New("Incorrect token format")
}
