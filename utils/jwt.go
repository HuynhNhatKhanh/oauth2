package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("khanh-secret-key")

type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

func CreateTokens(userID string) (string, string, error) {
	accessToken, err := GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	// refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, "", nil
}

func GenerateAccessToken(userID string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Thời gian hết hạn của token

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, nil
	} else {
		return "", err
	}
}

// func GenerateRefreshToken(userID string) (string, error) {
// 	// Tạo token với chuỗi ngẫu nhiên
// 	refreshToken
// 	// Lưu token vào database
// 	// Trả về token

// 	return refreshToken, nil
// }
