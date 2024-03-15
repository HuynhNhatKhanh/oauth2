package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashString(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func CompareStringHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
