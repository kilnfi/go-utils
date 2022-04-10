package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func New(cfg *Config) (*logrus.Logger, error) {
	var formatter logrus.Formatter

	switch cfg.Format {
	case "text":
		formatter = &logrus.TextFormatter{}
	case "json", "":
		formatter = &logrus.JSONFormatter{}
	default:
		return nil, fmt.Errorf("invalid log encoding format %q", cfg.Format)
	}

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	return &logrus.Logger{
		Formatter: formatter,
		Level:     level,
	}, nil
}
