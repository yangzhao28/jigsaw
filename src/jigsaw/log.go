package jigsaw

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyMsg:   "@message",
		},
	})
}
