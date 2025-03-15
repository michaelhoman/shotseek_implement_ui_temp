package mailer

type Client interface {
	Send(templateFile, first_name, last_name, email string, data any, isSandbox bool) error
}
