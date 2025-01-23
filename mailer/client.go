package mailer

import (
	"context"
	"errors"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sstp105/bangumi-component/utils"
	"time"
)

type MailgunConfig struct {
	ApiKey                 string   `json:"api_key" firestore:"api_key"`
	Domain                 string   `json:"domain" firestore:"domain"`
	From                   string   `json:"from" firestore:"from"`
	NotificationRecipients []string `json:"notification_recipients" firestore:"notification_recipients"`
}

// Client is a wrapper around the Firestore client and providing additional logging
type Client struct {
	cfg     *MailgunConfig
	mailgun *mailgun.MailgunImpl
}

// New creates a new mailer client.
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

// NotifyRecipients attempts to send the message given the configured recipients. The following are returned:
// - status: message deliver status
// - id: message id from mailgun
// - error: any errors if unable to send the message
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
	date := time.Now().Format(utils.YYYYMMDDDateFormatter)

	subject := fmt.Sprintf("%s Failed - %s", serviceName, date)
	message := fmt.Sprintf("Please visit Google Cloud > Logs for more details. Error: %s", e.Error())

	_, _, err := c.NotifyRecipients(ctx, subject, message)
	if err != nil {
		fmt.Printf("error sending emails to recipients:%s", err.Error())
	}
}
