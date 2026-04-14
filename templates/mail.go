package templates

import (
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {

	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	smtpHost := os.Getenv("HOST")
	smtpPort := os.Getenv("PORT")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		[]byte(msg),
	)
}
