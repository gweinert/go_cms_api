package services

import (
	"os"
	"strings"

	"gopkg.in/mailgun/mailgun-go.v1"
)

type Mail struct {
	ToEmail string
	Name    string
	Subject string
	Body    string
	Html    string
}

func SendSimpleMessage(mail *Mail) (string, error) {
	domain := os.Getenv("MG_DOMAIN")
	apiKey := os.Getenv("MG_API_KEY")
	publicApiKey := os.Getenv("MG_PUBLIC_API_KEY")

	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	from := strings.Join([]string{mail.Name, "<mailgun@YOUR_DOMAIN_NAME>"}, " ")
	m := mg.NewMessage(
		// "Excited User <mailgun@YOUR_DOMAIN_NAME>",
		from,
		mail.Subject,
		mail.Body,
		mail.ToEmail,
	)

	if mail.Html != "" {
		m.SetHtml(mail.Html)
	}

	_, id, err := mg.Send(m)
	return id, err
}
