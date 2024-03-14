package utils

import (
	"fmt"
	"math/rand"
	"time"

	gomail "gopkg.in/gomail.v2"
)

var ()

func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}

func SendOTP(email, otp string) error {

	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "OTP Verification")
	msg.SetBody("text/html", "Your OTP code is: "+otp)

	mail := gomail.NewDialer(host, port, from, key)

	if err := mail.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
