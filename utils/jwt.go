package utils

import (
	"fmt"
	"os"
	"time"
	"user_login/models"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(tokenType string, expiration time.Duration, user models.User) (string, error) {
	ecdsaPrivateKey, err := os.ReadFile("private.pem")
	if err != nil {
		return "", err
	}
	var privateKey interface{}
	if tokenType == "access" {
		er := error(nil)
		privateKey, er = jwt.ParseECPrivateKeyFromPEM(ecdsaPrivateKey)
		if er != nil {
			return "", er
		}
	} else if tokenType == "refresh" {
		er := error(nil)
		privateKey, er = jwt.ParseECPrivateKeyFromPEM(ecdsaPrivateKey)
		if er != nil {
			return "", er
		}
	}

	token := jwt.New(jwt.SigningMethodES256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(expiration).Unix()
	claims["iat"] = time.Now().Unix()
	claims["username"] = user.Username
	claims["email"] = user.Email

	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	// Parse public key
	ecdsaPublickey, err := os.ReadFile("public.pem")
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseECPublicKeyFromPEM(ecdsaPublickey)
	if err != nil {
		return nil, err
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

// func GenerateAccessToken(userID string) (string, error) {
// 	expirationTime := time.Now().Add(15 * time.Minute) // Thời gian hết hạn của token

// 	claims := &Claims{
// 		UserID: userID,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 			IssuedAt:  time.Now().Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	accessToken, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	return accessToken, nil
// }

// func ParseAccessToken(tokenString string) (string, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
// 		return claims.UserID, nil
// 	} else {
// 		return "", err
// 	}
// }

// func GenerateRefreshToken(userID string) (string, error) {
// 	// Tạo token với chuỗi ngẫu nhiên
// 	refreshToken
// 	// Lưu token vào database
// 	// Trả về token

// 	return refreshToken, nil
// }
