package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString hashes a string
func HashString(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

// CompareStringHash compares a hashed string with a string
func CompareStringHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
