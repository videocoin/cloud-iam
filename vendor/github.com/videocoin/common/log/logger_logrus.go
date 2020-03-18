package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogrus makes a new Interface backed by a logrus logger
func NewLogrus(level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stderr
	log.Level = level
	return log
}
