package middlewares

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)

func resendClient() (*resend.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	resend_api_key := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(resend_api_key)
	return client, nil
}

func SendEmail(email string, name string) error {
	client, clientErr := resendClient()
	if clientErr != nil {
		return clientErr
	}
	var html = "<h1>Hi " + name + "</h1>" + "<br> <h2>Verify your email</h2>"
	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{email},
		Subject: "Verification of email",
		Html:    html,
	}
	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}
	return nil
}
