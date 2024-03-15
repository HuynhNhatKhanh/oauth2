package utils

import (
	"fmt"
	"os"
	"time"
	"user_login/models"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a new JWT token
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

// ParseToken parses a JWT token
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
