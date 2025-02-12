package mailer

import (
	"context"
	"errors"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

type MailgunConfig struct {
	ApiKey                 string   `json:"api_key" firestore:"api_key"`
	Domain                 string   `json:"domain" firestore:"domain"`
	From                   string   `json:"from" firestore:"from"`
	NotificationRecipients []string `json:"notification_recipients" firestore:"notification_recipients"`
}

type Client struct {
	cfg     *MailgunConfig
	mailgun *mailgun.MailgunImpl
}

func New(cfg *MailgunConfig) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("mailgun config cannot be empty")
	}

	m := mailgun.NewMailgun(cfg.Domain, cfg.ApiKey)

	return &Client{
		cfg:     cfg,
		mailgun: m,
	}, nil
}

func (c *Client) NotifyRecipients(ctx context.Context, subject string, message string) (string, string, error) {
	m := mailgun.NewMessage(
		c.cfg.From,
		subject,
		message,
		c.cfg.NotificationRecipients...,
	)

	status, id, err := c.mailgun.Send(ctx, m)
	if err != nil {
		return "", "", err
	}

	return status, id, nil
}

func (c *Client) SendGoogleCloudFailureEmails(ctx context.Context, serviceName string, e error) {
	date := time.Now().Format("2006-01-02")

	subject := fmt.Sprintf("%s Failed - %s", serviceName, date)
	message := fmt.Sprintf("Please visit Google Cloud > Logs for more details. Error: %s", e.Error())

	_, _, err := c.NotifyRecipients(ctx, subject, message)
	if err != nil {
		fmt.Printf("error sending emails to recipients:%s", err.Error())
	}
}
