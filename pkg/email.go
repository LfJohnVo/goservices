package pkg

import (
	"fiber-boilerplate/config"
	"log"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

// SendEmail sends an email
func SendEmail(to []string, subject string, body string, html bool) error {
	m := gomail.NewMessage()
	// Set E-Mail sender
	m.SetHeader("From", config.GetEnvValue("MAIL_USERNAME"))

	// Set E-Mail receivers
	//setheader to array
	m.SetHeader("To", to...)

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	// Set E-Mail body. You can set plain text or html with text/html
	if html {
		m.SetBody("text/html", body)
	} else {
		m.SetBody("text/plain", body)
	}

	// Set E-Mail attachment
	// if file != "" {
	// 	m.Attach(file)
	// }

	// Settings for SMTP server
	port, _ := strconv.Atoi(config.GetEnvValue("MAIL_PORT"))
	d := gomail.NewDialer(config.GetEnvValue("MAIL_HOST"), port, config.GetEnvValue("MAIL_USERNAME"), config.GetEnvValue("MAIL_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		panic(err)
	}

	return nil
}
