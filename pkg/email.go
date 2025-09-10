package pkg

import (
	"github.com/horlerdipo/todo-golang/env"
	gomail "gopkg.in/mail.v2"
	"log"
	"strings"
)

type SendEmailConfig struct {
	Recipients  []string
	Subject     string
	Content     string
	ContentType string
}

func SendEmail(emailConfig SendEmailConfig) error {
	message := gomail.NewMessage()

	emailTo := strings.Join(emailConfig.Recipients, ",")

	message.SetHeader("From", env.FetchString("MAIL_FROM_ADDRESS"))
	message.SetHeader("To", emailTo)
	message.SetHeader("Subject", emailConfig.Subject)

	// Set email body
	message.SetBody(emailConfig.ContentType, emailConfig.Content)

	dialer := gomail.NewDialer(env.FetchString("MAIL_HOST"), env.FetchInt("MAIL_PORT"), env.FetchString("MAIL_USERNAME"), env.FetchString("MAIL_PASSWORD"))

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Error:", err)
		return err
	}

	log.Println("Email sent successfully to recipients:", emailConfig.Recipients)
	return nil
}
