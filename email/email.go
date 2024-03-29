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
	"github.com/alock/lotto-alert/util"
)

const (
	// smtp server configuration.
	smtpHost = "smtp.mail.me.com"
	smtpPort = "587"
)

func SendEmail(toEmail, message, dotenv string) error {
	err := godotenv.Load(dotenv)
	if err != nil {
		fmt.Println(err)
		return errors.New("error loading .env file")
	}
	appToken, ok := os.LookupEnv("APP_SPECIFIC")
	if !ok {
		return errors.New("no valid email app token")
	}
	// apple mail server docs - https://support.apple.com/en-us/HT202304
	auth := smtp.PlainAuth("", config.EmailStruct.From, appToken, smtpHost)
	// stackoverflow post helping figure out how to change the sender
	// https://stackoverflow.com/questions/71948786/how-to-use-smtp-with-apple-icloud-custom-domain
	rfc822style := fmt.Sprintf("From: %s\nTo: %s\nSubject: Wildlife Works Winner\n\n%s", config.EmailStruct.FromOverride, toEmail, message)
	var body bytes.Buffer
	body.Write([]byte(rfc822style))
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, config.EmailStruct.From, []string{toEmail}, body.Bytes())
	if err != nil {
		return err
	}
	fmt.Println("email sent")
	return nil
}

func GetMessage(today time.Time, winningNumber int, prizeInfo config.PrizeInfo) string {
	msg := fmt.Sprintf("Congrats! On %s the PICK 3 Evening Number was %s and you won $%v", today.Format("January 2"), util.PadLottoInt(winningNumber), prizeInfo.Amount)
	if prizeInfo.Reason != "" {
		msg = msg + fmt.Sprintf(" because it is %s", prizeInfo.Reason)
	}
	return msg + "."
}
