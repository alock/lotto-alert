package email

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/alock/lotto-alert/config"
)

const (
	// smtp server configuration.
	smtpHost = "smtp.mail.me.com"
	smtpPort = "587"
)

func SendEmail(email, message, dotenv string) error {
	err := godotenv.Load(dotenv)
	if err != nil {
		return errors.New("error loading .env file")
	}
	// apple mail server docs - https://support.apple.com/en-us/HT202304
	to := []string{email}
	auth := smtp.PlainAuth("", config.EmailStruct.From, os.Getenv("APP_SPECIFIC"), smtpHost)
	// stackoverflow post helping figure out how to change the sender
	// https://stackoverflow.com/questions/71948786/how-to-use-smtp-with-apple-icloud-custom-domain
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n"
	body.Write([]byte(fmt.Sprintf("To: %s\nFrom: %s\nSubject: Wildlife Works Winner \n%s%s", to[0], config.EmailStruct.FromOverride, mimeHeaders, message)))
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, config.EmailStruct.From, to, body.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("email sent")
	return nil
}

func GetMessage(today time.Time, winningNumber int, prizeInfo config.PrizeInfo) string {
	msg := fmt.Sprintf("Congrats! On %s the PICK 3 Evening Number was %v and you won $%v", today.Format("January 2"), winningNumber, prizeInfo.Amount)
	if prizeInfo.Reason != "" {
		msg = msg + fmt.Sprintf(" because it is %s", prizeInfo.Reason)
	}
	return msg + "."
}
