package interfaces

import "github.com/sirupsen/logrus"

type Loggable interface {
	Logger() *logrus.Logger
	SetLogger(*logrus.Logger)
}
