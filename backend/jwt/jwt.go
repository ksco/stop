package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func GenJwt(id string, secret []byte, days time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(24 * time.Hour * days).Unix(),
	})

	t, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return t, nil
}

func ParseJwt(tokenStr string, secret []byte) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		v, ok := claims["id"]
		if !ok {
			return "", errors.New("parse token err")
		}
		return v.(string), nil
	}

	return "", errors.New("parse token err")
}
