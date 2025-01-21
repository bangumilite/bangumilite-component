package handler

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-component/bangumi"
	"github.com/sstp105/bangumi-component/fs"
	"github.com/sstp105/bangumi-component/mailer"
	"github.com/sstp105/bangumi-component/model"
	"os"
	"time"
)

const (
	GoogleCloudFieldKeyLevel = "severity"
	GoogleCloudFieldKeyMsg   = "message"

	DateFormatter = "2006-01-02"
)

type GCPHandler struct {
	serviceName   string
	environment   model.Environment
	logger        *logrus.Logger
	bangumiClient *bangumi.Client
	fsClient      *fs.Client
	mailerClient  *mailer.Client
}

type Option func(h *GCPHandler)

// New creates a new GCPHandler that initialize all the necessary clients for Google cloud serverless functions
func New(options ...Option) (*GCPHandler, error) {

	ctx := context.Background()

	// environment setup
	var env model.Environment
	if os.Getenv(model.RunningEnvironment) == string(model.Production) {
		env = model.Production
	} else {
		env = model.Local
	}

	// logger setup - map log filed key to google cloud compatible key
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   GoogleCloudFieldKeyMsg,
			logrus.FieldKeyLevel: GoogleCloudFieldKeyLevel,
		},
	})

	// initialize bangumi client
	bangumiClient := bangumi.NewClient()

	// initialize firestore client
	fsClient, err := fs.New(ctx)
	if err != nil {
		return nil, err
	}

	// retrieve Mailgun credentials
	mailgunCfg, err := fsClient.GetMailgunConfig(ctx)
	if err != nil {
		return nil, err
	}

	// initialize Mailgun client for sending emails
	mailerClient, err := mailer.New(mailgunCfg)
	if err != nil {
		return nil, err
	}

	h := &GCPHandler{
		environment:   env,
		logger:        logger,
		bangumiClient: bangumiClient,
		fsClient:      fsClient,
		mailerClient:  mailerClient,
	}

	// apply caller passed custom options
	for _, opt := range options {
		opt(h)
	}

	return h, nil
}

func (h *GCPHandler) Sync(ctx context.Context, processFunc func(ctx context.Context) error) error {
	start := time.Now()

	h.logger.Infof("%s triggered, env: %s", h.serviceName, h.environment)

	err := processFunc(ctx)

	if err != nil {
		h.sendEmails(ctx, err)
		return err
	}

	duration := time.Since(start)
	h.logger.Infof("%s completed, total execuation time:%s", h.serviceName, duration)

	return nil
}

func WithServiceName(name string) Option {
	return func(h *GCPHandler) {
		h.serviceName = name
	}
}

func WithLogLevel(level logrus.Level) Option {
	return func(h *GCPHandler) {
		h.logger.SetLevel(level)
	}
}

func (h *GCPHandler) sendEmails(ctx context.Context, e error) {
	date := time.Now().Format(DateFormatter)

	subject := fmt.Sprintf("%s Failed - %s", h.serviceName, date)
	message := fmt.Sprintf("Please visit Google Cloud > Logs for more details. Error: %s", e.Error())

	_, _, err := h.mailerClient.NotifyRecipients(ctx, subject, message)
	if err != nil {
		h.logger.Errorf("error sending emails to recipients, error:%s", err)

		// fallthrough if unable to send emails to recipients
	}
}
