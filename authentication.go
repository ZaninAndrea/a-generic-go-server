package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type AccessToken struct {
	UserID string
}

func parseBearer(authorizationHeader []string) (*AccessToken, error) {
	if len(authorizationHeader) == 0 {
		return nil, fmt.Errorf("You need to pass an authentication token in the Authorization header")
	} else if !strings.HasPrefix(authorizationHeader[0], "Bearer ") {
		return nil, fmt.Errorf(`Authorization token should have the "Bearer " prefix`)
	}
	token := authorizationHeader[0][7:]
	parsedToken, err := VerifyToken(token)

	if err != nil {
		return nil, fmt.Errorf("Malformed or invalid token")
	}

	return parsedToken, nil
}

func VerifyToken(tokenString string) (*AccessToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		return nil, err
	}

	return &AccessToken{userID}, nil
}

func GenerateToken(id string) string {
	atClaims := jwt.MapClaims{"userId": id}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		panic(err)
	}

	return token
}
