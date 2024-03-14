package utils

import (
	"fmt"
	"math/rand"

	gomail "gopkg.in/gomail.v2"
)

var (
	key  = "xsmtpsib-2c6d85049fa4c99f949f4e54fa40cff35b33dc1260c3d65c9732b4ba7ae8e56a-8d2cyLZOGVm43WgX"
	from = "khanhhuynh28082000@gmail.com"
	host = "smtp-relay.brevo.com"
	port = 587
)

func GenerateOTP() string {
	// rand.Seed(time.Now().UnixNano())
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
		fmt.Println(err)
		return err
	}

	return nil
}
