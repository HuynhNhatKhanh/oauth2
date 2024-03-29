package utils

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

// GenerateOTP generates a random 6 digit number
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}

// SendLinkOrOTP sends link or to the user's email
func SendLinkOrOTP(email, options string, otp string, username string) error {

	hostLink := os.Getenv("HOST")

	from := os.Getenv("MAIL_FROM_ADDRESS")
	pass := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	to := email
	address := host + ":" + port

	msg := ""

	if options == "link" {
		msg = "From: " + from + "\n" +
			"To: " + to + "\n" +
			"Subject: Email Verification\n\n" +
			"Click the link to verify your email: " + hostLink + "/verify?username=" + username + "&email=" + otp

	} else if options == "otp" {
		msg = "From: " + from + "\n" +
			"To: " + to + "\n" +
			"Subject: OTP Verification\n\n" +
			"Your OTP code is: " + otp
	}

	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(address, auth, from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
