package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var (
	JWT_SECRET = os.Getenv("JWT_SECRET")
)

func Parse(credential string) (*User, error) {
	token, err := jwt.Parse(credential, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SECRET), nil
	})

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims jwt error")
	}

	u := &User{}

	sub := fmt.Sprintf("%v", claims["sub"])

	if !ok {
		return nil, errors.New("claims jwt error, sub type not a string")
	}

	err = json.Unmarshal([]byte(sub), u)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func CheckDefaultAuth(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token not exist")
	}
	u, err := Parse(token)
	if err != nil {
		return "", err
	}
	return u.Id, nil
}
