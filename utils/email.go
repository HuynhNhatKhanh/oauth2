package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}

func SendOTP(email, otp string) error {

	from := os.Getenv("MAIL_FROM_ADDRESS")
	pass := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	to := email

	address := host + ":" + port
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OTP Verification\n\n" +
		"Your OTP code is: " + otp

	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(address, auth, from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
