package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"text/template"
	"time"

	"github.com/sandoxlabs99/gopher_social/internal/store"
	_ "github.com/sandoxlabs99/gopher_social/internal/store"

	"github.com/resend/resend-go/v3"
	"go.uber.org/zap"
)

type ResendClient struct {
	fromEmail string
	apiKey    string
	client    *resend.Client
	logger    *zap.SugaredLogger
}

func NewResendClient(apiKey, fromEmail string, logger *zap.SugaredLogger) (*ResendClient, error) {
	if apiKey == "" {
		return &ResendClient{}, errors.New("api key is required")
	}

	client := resend.NewClient(apiKey)

	return &ResendClient{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
		logger:    logger,
	}, nil
}

func (m *ResendClient) Send(templateFile, username, email string, data any, isSandbox bool) error {
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Html:    body.String(),
		Subject: subject.String(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), store.QueryTimeoutDuration)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond) // use to test for email failure
	defer cancel()

	for i := range maxRetries {
		response, err := m.client.Emails.SendWithContext(ctx, params)
		if err != nil {
			errMsg := fmt.Sprintf("failed to send an email, attempt %d of %d", i+1, maxRetries)
			m.logger.Warnw(
				errMsg,
				"error", err.Error(),
			)

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		m.logger.Infof("Email sent with response id %v", response.Id)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts, error: %v", maxRetries, err)
}
