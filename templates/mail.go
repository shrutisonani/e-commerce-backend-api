package templates

import (
	"fmt"
	"net/smtp"
	"os"
)

// send email using SMTP server configured in environment variables
func SendEmail(to, subject, body string) error {

	from := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	smtpHost := os.Getenv("HOST")
	smtpPort := os.Getenv("PORT")

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
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

// generate email body for verification email
func GenerateVerificationEmail(token string) string {

	link := fmt.Sprintf("http://localhost:8080/verify-email?token=%s", token)

	return fmt.Sprintf(`
		<html>
		<body>
			<h2>Email Verification</h2>
			<p>Please verify your email within 15 minutes.</p>

			<a href="%s" style="
				display:inline-block;
				padding:10px 20px;
				background-color:#4CAF50;
				color:white;
				text-decoration:none;
				border-radius:5px;">
				Verify Email
			</a>

			<p>If button doesn't work, click this link:</p>
			<p>%s</p>
		</body>
		</html>
	`, link, link)
}
