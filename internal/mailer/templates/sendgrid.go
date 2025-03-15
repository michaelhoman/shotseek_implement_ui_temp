package mailer

import "github.com/sendgrid/sendgrid-go"

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}
