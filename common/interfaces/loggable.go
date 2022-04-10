package interfaces

import "github.com/sirupsen/logrus"

type Loggable interface {
	Logger() logrus.FieldLogger
	SetLogger(logrus.FieldLogger)
}
