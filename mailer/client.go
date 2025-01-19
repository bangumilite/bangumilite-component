package mailer

import (
	"context"
	"errors"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
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
	logger  *logrus.Logger
	mailgun *mailgun.MailgunImpl
}

// New creates a new mailer client.
func New(cfg *MailgunConfig, logger *logrus.Logger) (*Client, error) {
	if cfg == nil {
		return nil, errors.New("mailgun config cannot be empty")
	}

	m := mailgun.NewMailgun(cfg.Domain, cfg.ApiKey)

	return &Client{
		cfg:     cfg,
		logger:  logger,
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
		c.logger.Errorf("failed to notify recipients, error:%s", err.Error())
		return "", "", err
	}

	return status, id, nil
}
