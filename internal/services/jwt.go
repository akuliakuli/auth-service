package services

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	UserID string `json:"userId"`
	IP     string `json:"ip"`
	jwt.RegisteredClaims
}

// Генерация JWT токена
func GenerateAccessToken(userID, ip string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		UserID: userID,
		IP:     ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(jwtKey)
}

// Валидация JWT токена
func ValidateAccessToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// Генерация случайного Refresh токена
func GenerateRefreshToken() (string, string, error) {
	refreshToken := make([]byte, 64)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", "", err
	}

	hashedRefreshToken := sha512.Sum512(refreshToken)
	encodedToken := base64.StdEncoding.EncodeToString(refreshToken)
	return encodedToken, string(hashedRefreshToken[:]), nil
}
