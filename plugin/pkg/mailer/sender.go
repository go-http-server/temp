package mailer

import (
	"fmt"
	"html/template"

	"github.com/wneessen/go-mail"
)

type UserReceive struct {
	Username     string   `json:"username"`
	EmailAddress string   `json:"email_address"`
	Code         string   `json:"code"`
	Fullname     string   `json:"full_name"`
	AttachFiles  []string `json:"attach_files"`
}

type EmailSender interface {
	SendWithTemplate(subject, pathTemplate string, to UserReceive) error
}

type GmailSender struct {
	name             string
	emailAddress     string
	emailAppPassword string
}

func NewGmailSender(name, emailAddress, emailAppPassword string) EmailSender {
	return &GmailSender{
		name:             name,
		emailAddress:     emailAddress,
		emailAppPassword: emailAppPassword,
	}
}

// WARN: pathTemplate is exactly is root project, example: ./templates/verify_account.html. If use test, U must be use absolute path.
func (sender GmailSender) SendWithTemplate(subject, pathTemplate string, receiver UserReceive) error {
	htmlTemplate, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return fmt.Errorf("error with path template file: %s", err)
	}
	message := mail.NewMsg()

	if err := message.EnvelopeFrom(sender.emailAddress); err != nil {
		return fmt.Errorf("failed to set ENVELOPE FROM address: %s", err)
	}
	if err := message.FromFormat(sender.name, sender.emailAddress); err != nil {
		return fmt.Errorf("failed to set formatted FROM address: %s", err)
	}
	if err := message.AddToFormat(receiver.Fullname, receiver.EmailAddress); err != nil {
		return fmt.Errorf("failed to set formatted TO address: %s", err)
	}

	message.SetMessageID()
	message.SetDate()

	message.Subject(subject)
	if err := message.AddAlternativeHTMLTemplate(htmlTemplate, receiver); err != nil {
		return fmt.Errorf("failed to add HTML template to mail body: %s", err)
	}

	for _, filePath := range receiver.AttachFiles {
		message.AttachFile(filePath, mail.WithFileDescription("From Go with love"))
	}

	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(sender.emailAddress),
		mail.WithPassword(sender.emailAppPassword),
	)
	if err != nil {
		return fmt.Errorf("cannot create client mailer: %s", err)
	}

	if err = client.DialAndSend(message); err != nil {
		return fmt.Errorf("cannot delivery mail: %s", err)
	}
	return nil
}
