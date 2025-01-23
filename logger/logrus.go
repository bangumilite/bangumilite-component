package logger

import "github.com/sirupsen/logrus"

const (
	GoogleCloudFieldKeyLevel = "severity"
	GoogleCloudFieldKeyMsg   = "message"
)

func NewLogrus(level logrus.Level) *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   GoogleCloudFieldKeyMsg,
			logrus.FieldKeyLevel: GoogleCloudFieldKeyLevel,
		},
	})

	logger.SetLevel(level)

	return logger
}
